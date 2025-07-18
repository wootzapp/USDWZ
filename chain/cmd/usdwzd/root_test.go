package main

import "testing"

func TestNewRootCmd(t *testing.T) {
	cmd := newRootCmd()
	if cmd.Use != "usdwzd" {
		t.Fatalf("expected use 'usdwzd' got %s", cmd.Use)
	}
}

func TestDefaultNodeHome(t *testing.T) {
	if DefaultNodeHome == "" {
		t.Fatal("default node home should not be empty")
	}
}
