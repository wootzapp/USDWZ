package types

import "testing"

func TestModuleName(t *testing.T) {
	if ModuleName != "yield" {
		t.Fatalf("expected yield got %s", ModuleName)
	}
}
