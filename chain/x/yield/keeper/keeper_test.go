package keeper

import (
	"fmt"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stablecoinkeeper "github.com/example/usdwz/chain/x/stablecoin/keeper"
	"github.com/example/usdwz/chain/x/yield/types"
)

func TestNewKeeper(t *testing.T) {
	key := storetypes.NewKVStoreKey("yield")
	k := NewKeeper(nil, key)
	if k.storeKey != key {
		t.Fatal("store key mismatch")
	}
}

func TestModuleName(t *testing.T) {
	if ModuleName() != types.ModuleName {
		t.Fatalf("expected %s got %s", types.ModuleName, ModuleName())
	}
}

func TestAccrueAndDistribute(t *testing.T) {
	key := storetypes.NewKVStoreKey("yield")
	k := NewKeeper(nil, key)
	ctx := stablecoinkeeper.TestContext(key)

	k.Accrue(ctx, sdk.NewInt(5))
	if !k.Accrued(ctx).Equal(sdk.NewInt(5)) {
		t.Fatal("accrue failed")
	}
	k.Accrue(ctx, sdk.NewInt(3))
	if !k.Accrued(ctx).Equal(sdk.NewInt(8)) {
		t.Fatal("second accrue failed")
	}
	dist := k.Distribute(ctx)
	if !dist.Equal(sdk.NewInt(8)) || !k.Accrued(ctx).IsZero() {
		t.Fatal("distribute incorrect")
	}
}

func TestYieldEdgeCases(t *testing.T) {
	key := storetypes.NewKVStoreKey("yield")
	k := NewKeeper(nil, key)
	ctx := stablecoinkeeper.TestContext(key)

	cases := []sdk.Int{sdk.ZeroInt(), sdk.NewInt(-1), sdk.NewInt(100)}
	for i, amt := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			k.Accrue(ctx, amt)
			_ = k.Distribute(ctx)
		})
	}
}
