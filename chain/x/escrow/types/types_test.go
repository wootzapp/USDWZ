package types

import "testing"

func TestModuleName(t *testing.T) {
	if ModuleName != "escrow" {
		t.Fatalf("expected escrow got %s", ModuleName)
	}
}
