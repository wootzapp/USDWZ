package types

import "testing"

func TestModuleName(t *testing.T) {
	if ModuleName != "audit" {
		t.Fatalf("expected audit got %s", ModuleName)
	}
}

func TestAttestation(t *testing.T) {
	a := Attestation{Timestamp: 1, Hash: "abc"}
	if a.Timestamp != 1 || a.Hash != "abc" {
		t.Fatal("attestation fields not assigned")
	}
}
