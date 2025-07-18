package types

import "testing"

func TestModuleName(t *testing.T) {
	if ModuleName != "validator" {
		t.Fatalf("expected validator got %s", ModuleName)
	}
}
