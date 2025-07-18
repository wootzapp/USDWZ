//go:build integration

package integration_test

import (
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/example/usdwz/chain/app"
	escrowtypes "github.com/example/usdwz/chain/x/escrow/types"
	stabletypes "github.com/example/usdwz/chain/x/stablecoin/types"
	validatortypes "github.com/example/usdwz/chain/x/validator/types"
	yieldtypes "github.com/example/usdwz/chain/x/yield/types"
)

func contextForApp(a *app.UsdWzApp) sdk.Context {
	db := dbm.NewMemDB()
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
	return sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())
}

func TestStablecoinE2E(t *testing.T) {
	a := app.New(log.NewNopLogger())
	ctx := contextForApp(a)

	a.StablecoinKeeper.DepositCollateral(ctx, sdk.NewInt(100))
	if !a.StablecoinKeeper.GetCollateral(ctx).Equal(sdk.NewInt(100)) {
		t.Fatal("collateral expected")
	}

	if err := a.StablecoinKeeper.Mint(ctx, sdk.NewInt(50)); err != nil {
		t.Fatalf("mint failed: %v", err)
	}

	if err := a.StablecoinKeeper.RedeemCollateral(ctx, sdk.NewInt(60)); err == nil {
		t.Fatal("redeem should fail")
	}

	if err := a.StablecoinKeeper.Burn(ctx, sdk.NewInt(20)); err != nil {
		t.Fatalf("burn failed: %v", err)
	}

	if !a.StablecoinKeeper.GetSupply(ctx).Equal(sdk.NewInt(30)) {
		t.Fatal("unexpected supply")
	}
}
