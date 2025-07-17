package keeper

import (
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
