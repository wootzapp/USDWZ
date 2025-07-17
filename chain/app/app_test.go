package app_test

import (
	"testing"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/example/usdwz/chain/app"
	escrowtypes "github.com/example/usdwz/chain/x/escrow/types"
	stablecointypes "github.com/example/usdwz/chain/x/stablecoin/types"
)

func TestNewApp(t *testing.T) {
	a := app.New(log.NewNopLogger())
	if a.Name() != "usdWz" {
		t.Fatalf("unexpected app name: %s", a.Name())
	}

	if a.BaseApp == nil {
		t.Fatal("base app not initialized")
	}
}

func TestModuleManager(t *testing.T) {
	a := app.New(log.NewNopLogger())
	if a.ModuleManager() == nil {
		t.Fatal("module manager should not be nil")
	}
}

func TestKVStoreKey(t *testing.T) {
	a := app.New(log.NewNopLogger())
	k1 := a.KVStoreKey(stablecointypes.StoreKey)
	k2 := a.KVStoreKey(escrowtypes.StoreKey)
	if k1 == nil || k2 == nil {
		t.Fatal("store keys should be created")
	}
	if k1 == k2 {
		t.Fatal("store keys should be distinct")
	}
}
