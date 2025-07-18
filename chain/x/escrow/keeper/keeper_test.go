package keeper

import (
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/example/usdwz/chain/x/escrow/types"
	stablecoinkeeper "github.com/example/usdwz/chain/x/stablecoin/keeper"
)

func TestNewKeeper(t *testing.T) {
	key := storetypes.NewKVStoreKey("escrow")
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

func TestEscrowLifecycle(t *testing.T) {
	key := storetypes.NewKVStoreKey("escrow")
	k := NewKeeper(nil, key)
	ctx := stablecoinkeeper.TestContext(key) // reuse helper

	k.CreateEscrow(ctx, "e1", sdk.NewInt(10))
	e, ok := k.GetEscrow(ctx, "e1")
	if !ok || !e.Amount.Equal(sdk.NewInt(10)) {
		t.Fatal("escrow not stored")
	}
	if e.Completed {
		t.Fatal("should not be completed")
	}
	if err := k.FinalizeEscrow(ctx, "e1"); err != nil {
		t.Fatalf("finalize failed: %v", err)
	}
	e, _ = k.GetEscrow(ctx, "e1")
	if !e.Completed {
		t.Fatal("escrow not finalized")
	}
	k.DeleteEscrow(ctx, "e1")
	if _, ok := k.GetEscrow(ctx, "e1"); ok {
		t.Fatal("escrow should be deleted")
	}
}

func TestEscrowList(t *testing.T) {
	key := storetypes.NewKVStoreKey("escrow")
	k := NewKeeper(nil, key)
	ctx := stablecoinkeeper.TestContext(key)

	ids := []string{"a", "b", "c"}
	for i, id := range ids {
		k.CreateEscrow(ctx, id, sdk.NewInt(int64(i)))
	}
	esc := k.ListEscrows(ctx)
	if len(esc) != len(ids) {
		t.Fatalf("expected %d got %d", len(ids), len(esc))
	}
}

func TestEscrowEdgeCases(t *testing.T) {
	key := storetypes.NewKVStoreKey("escrow")
	k := NewKeeper(nil, key)
	ctx := stablecoinkeeper.TestContext(key)

	cases := []struct {
		id  string
		amt sdk.Int
	}{
		{"zero", sdk.ZeroInt()},
		{"negative", sdk.NewInt(-1)},
		{"big", sdk.NewInt(1000)},
		{"dup", sdk.OneInt()},
		{"dup", sdk.NewInt(2)},
		{"spaces id", sdk.NewInt(3)},
		{"unicode☃", sdk.NewInt(4)},
	}

	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			k.CreateEscrow(ctx, tc.id, tc.amt)
			esc, _ := k.GetEscrow(ctx, tc.id)
			if !esc.Amount.Equal(tc.amt) {
				t.Fatalf("amt mismatch %s", tc.id)
			}
		})
	}
}
