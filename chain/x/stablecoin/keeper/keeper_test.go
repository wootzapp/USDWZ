package keeper

import (
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/example/usdwz/chain/x/stablecoin/types"
)

func TestNewKeeper(t *testing.T) {
	key := storetypes.NewKVStoreKey("test")
	k := NewKeeper(nil, key)
	if k.storeKey != key {
		t.Fatal("store key not assigned")
	}
}

func TestModuleName(t *testing.T) {
	if ModuleName() != types.ModuleName {
		t.Fatalf("expected %s got %s", types.ModuleName, ModuleName())
	}
}

func TestMintAndCollateral(t *testing.T) {
	key := storetypes.NewKVStoreKey("sc")
	ctx := TestContext(key)
	k := NewKeeper(nil, key)

	k.DepositCollateral(ctx, sdk.NewInt(100))
	if !k.GetCollateral(ctx).Equal(sdk.NewInt(100)) {
		t.Fatal("collateral not recorded")
	}

	if err := k.Mint(ctx, sdk.NewInt(110)); err == nil {
		t.Fatal("should fail when collateral insufficient")
	}

	if err := k.Mint(ctx, sdk.NewInt(50)); err != nil {
		t.Fatalf("mint failed: %v", err)
	}
	if !k.GetSupply(ctx).Equal(sdk.NewInt(50)) {
		t.Fatal("supply mismatch after mint")
	}

	if err := k.RedeemCollateral(ctx, sdk.NewInt(60)); err == nil {
		t.Fatal("redeem should fail")
	}

	if err := k.RedeemCollateral(ctx, sdk.NewInt(25)); err != nil {
		t.Fatalf("redeem failed: %v", err)
	}
	if !k.GetCollateral(ctx).Equal(sdk.NewInt(75)) {
		t.Fatal("collateral not updated")
	}

	if err := k.Burn(ctx, sdk.NewInt(50)); err != nil {
		t.Fatalf("burn failed: %v", err)
	}
	if !k.GetSupply(ctx).Equal(sdk.ZeroInt()) {
		t.Fatal("supply not zero after burn")
	}
}
