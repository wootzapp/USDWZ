package escrow_test

import (
	"testing"

	"github.com/example/usdwz/chain/x/escrow"
	"github.com/example/usdwz/chain/x/escrow/keeper"
	"github.com/example/usdwz/chain/x/escrow/types"
)

func TestModuleName(t *testing.T) {
	if keeper.ModuleName() != types.ModuleName {
		t.Fatalf("expected %s got %s", types.ModuleName, keeper.ModuleName())
	}
}

func TestDefaultGenesis(t *testing.T) {
	var m escrow.AppModuleBasic
	if string(m.DefaultGenesis(nil)) != "{}" {
		t.Fatal("default genesis should be empty JSON")
	}
}
