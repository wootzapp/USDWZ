//go:build integration

package integration_test

import (
	"net"
	"net/http"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/example/usdwz/chain/app"
	"github.com/example/usdwz/chain/monitoring"
	escrowtypes "github.com/example/usdwz/chain/x/escrow/types"
	stabletypes "github.com/example/usdwz/chain/x/stablecoin/types"
	validatortypes "github.com/example/usdwz/chain/x/validator/types"
	yieldtypes "github.com/example/usdwz/chain/x/yield/types"
)

func contextForAppWithDB(a *app.UsdWzApp, db dbm.DB) (sdk.Context, store.CommitMultiStore) {
	ms := store.NewCommitMultiStore(db)
	keys := []string{
		stabletypes.StoreKey,
		escrowtypes.StoreKey,
		validatortypes.StoreKey,
		yieldtypes.StoreKey,
	}
	for _, name := range keys {
		ms.MountStoreWithDB(a.KVStoreKey(name), storetypes.StoreTypeIAVL, db)
	}
	_ = ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())
	return ctx, ms
}

func newAppCtx() (sdk.Context, *app.UsdWzApp) {
	a := app.New(log.NewNopLogger())
	ctx, _ := contextForAppWithDB(a, dbm.NewMemDB())
	return ctx, a
}

func TestSupplyZeroAtStart(t *testing.T) {
	ctx, a := newAppCtx()
	if !a.StablecoinKeeper.GetSupply(ctx).IsZero() {
		t.Fatal("supply should start at zero")
	}
}

func TestMintEqualToCollateral(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(10))
	if err := a.StablecoinKeeper.Mint(ctx, sdk.NewInt(10)); err != nil {
		t.Fatalf("mint failed: %v", err)
	}
	if !a.StablecoinKeeper.GetSupply(ctx).Equal(sdk.NewInt(10)) {
		t.Fatal("supply mismatch")
	}
}

func TestMintFailsWhenSupplyExceedsCollateral(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(5))
	if err := a.StablecoinKeeper.Mint(ctx, sdk.NewInt(6)); err == nil {
		t.Fatal("expected mint failure")
	}
}

func TestBurnFailsWithoutSupply(t *testing.T) {
	ctx, a := newAppCtx()
	if err := a.StablecoinKeeper.Burn(ctx, sdk.OneInt()); err == nil {
		t.Fatal("burn should fail with no supply")
	}
}

func TestBurnReducesSupplyProperly(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(20))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(10))
	_ = a.StablecoinKeeper.Burn(ctx, sdk.NewInt(5))
	if !a.StablecoinKeeper.GetSupply(ctx).Equal(sdk.NewInt(5)) {
		t.Fatal("burn did not reduce supply")
	}
}

func TestRedeemFullCollateralAfterBurn(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(20))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(10))
	_ = a.StablecoinKeeper.Burn(ctx, sdk.NewInt(10))
	if err := a.StablecoinKeeper.RedeemCollateral(ctx, sdk.NewInt(20)); err != nil {
		t.Fatalf("redeem failed: %v", err)
	}
}

func TestRedeemFailsWhenSupplyMoreThanCollateral(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(20))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(10))
	if err := a.StablecoinKeeper.RedeemCollateral(ctx, sdk.NewInt(15)); err == nil {
		t.Fatal("redeem should fail")
	}
}

func TestMultipleDepositsIncreaseCollateral(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(20))
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(30))
	if !a.StablecoinKeeper.GetCollateral(ctx).Equal(sdk.NewInt(50)) {
		t.Fatal("collateral total incorrect")
	}
}

func TestMultipleMintsIncreaseSupply(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(50))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(10))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(20))
	if !a.StablecoinKeeper.GetSupply(ctx).Equal(sdk.NewInt(30)) {
		t.Fatal("supply total incorrect")
	}
}

func TestMintAfterRedeemFailsWhenNotEnoughCollateral(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(30))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(30))
	if err := a.StablecoinKeeper.Mint(ctx, sdk.OneInt()); err == nil {
		t.Fatal("mint should fail")
	}
}

func TestMintAfterRedeemSuccessWhenEnoughCollateral(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(50))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(20))
	_ = a.StablecoinKeeper.RedeemCollateral(ctx, sdk.NewInt(15))
	if err := a.StablecoinKeeper.Mint(ctx, sdk.NewInt(10)); err != nil {
		t.Fatalf("mint failed: %v", err)
	}
}

func TestDepositNegativeAmount(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(-10))
	if !a.StablecoinKeeper.GetCollateral(ctx).Equal(sdk.NewInt(-10)) {
		t.Fatal("negative deposit not recorded")
	}
}

func TestMintZeroNoChange(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(10))
	if err := a.StablecoinKeeper.Mint(ctx, sdk.ZeroInt()); err != nil {
		t.Fatalf("mint zero failed: %v", err)
	}
	if !a.StablecoinKeeper.GetSupply(ctx).IsZero() {
		t.Fatal("supply should remain zero")
	}
}

func TestBurnZeroNoChange(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(10))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(5))
	if err := a.StablecoinKeeper.Burn(ctx, sdk.ZeroInt()); err != nil {
		t.Fatalf("burn zero failed: %v", err)
	}
	if !a.StablecoinKeeper.GetSupply(ctx).Equal(sdk.NewInt(5)) {
		t.Fatal("supply should remain 5")
	}
}

func TestCollateralNotLessThanSupplyInvariant(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(50))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(20))
	if a.StablecoinKeeper.GetCollateral(ctx).LT(a.StablecoinKeeper.GetSupply(ctx)) {
		t.Fatal("collateral less than supply")
	}
	_ = a.StablecoinKeeper.Burn(ctx, sdk.NewInt(10))
	if a.StablecoinKeeper.GetCollateral(ctx).LT(a.StablecoinKeeper.GetSupply(ctx)) {
		t.Fatal("collateral invariant broken after burn")
	}
}

func TestMetricsServerPortInUse(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer ln.Close()
	done := make(chan error, 1)
	go func() { done <- monitoring.StartServer(ln.Addr().String()) }()
	err = <-done
	if err == nil {
		t.Fatal("expected error when port in use")
	}
}

func TestNetworkDialFailure(t *testing.T) {
	_, err := http.Get("http://127.0.0.1:1/metrics")
	if err == nil {
		t.Fatal("expected dial error")
	}
}

func TestMetricsGaugeReflectsDeposit(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(20))
	if v := testutil.ToFloat64(monitoring.CollateralGaugeForTest()); v != 20 {
		t.Fatalf("expected gauge 20 got %v", v)
	}
}

func TestStablecoinOutOfPeg(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(10))
	if err := a.StablecoinKeeper.Mint(ctx, sdk.NewInt(10)); err != nil {
		t.Fatal("mint should succeed")
	}
	if err := a.StablecoinKeeper.Mint(ctx, sdk.OneInt()); err == nil {
		t.Fatal("expected mint failure when out of peg")
	}
}

func TestStablecoinMultipleOperations(t *testing.T) {
	ctx, a := newAppCtx()
	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(100))
	_ = a.StablecoinKeeper.Mint(ctx, sdk.NewInt(40))
	_ = a.StablecoinKeeper.Burn(ctx, sdk.NewInt(10))
	_ = a.StablecoinKeeper.RedeemCollateral(ctx, sdk.NewInt(20))
	if !a.StablecoinKeeper.GetSupply(ctx).Equal(sdk.NewInt(30)) {
		t.Fatal("unexpected final supply")
	}
	if !a.StablecoinKeeper.GetCollateral(ctx).Equal(sdk.NewInt(80)) {
		t.Fatal("unexpected final collateral")
	}
}
