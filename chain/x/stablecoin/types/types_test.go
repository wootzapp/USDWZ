package types

import "testing"

func TestModuleName(t *testing.T) {
	if ModuleName != "stablecoin" {
		t.Fatalf("expected stablecoin got %s", ModuleName)
	}
}
