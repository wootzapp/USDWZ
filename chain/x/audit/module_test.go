package audit

import "testing"

func TestModuleName(t *testing.T) {
	var m AppModuleBasic
	if m.Name() != "audit" {
		t.Fatalf("expected audit got %s", m.Name())
	}
}

func TestDefaultGenesis(t *testing.T) {
	var m AppModuleBasic
	if string(m.DefaultGenesis(nil)) != "{}" {
		t.Fatal("default genesis should be empty JSON")
	}
}
