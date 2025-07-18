package stablecoin_test

import (
	"testing"

	"github.com/example/usdwz/chain/x/stablecoin"
	"github.com/example/usdwz/chain/x/stablecoin/keeper"
	"github.com/example/usdwz/chain/x/stablecoin/types"
)

func TestModuleName(t *testing.T) {
	if keeper.ModuleName() != types.ModuleName {
		t.Fatalf("expected %s got %s", types.ModuleName, keeper.ModuleName())
	}
}

func TestDefaultGenesis(t *testing.T) {
	var m stablecoin.AppModuleBasic
	if string(m.DefaultGenesis(nil)) != "{}" {
		t.Fatal("default genesis should be empty JSON")
	}
}
