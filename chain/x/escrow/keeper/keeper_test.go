package keeper

import (
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/example/usdwz/chain/x/escrow/types"
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
