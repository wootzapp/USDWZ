//go:build integration

package integration_test

import (
	"testing"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/example/usdwz/chain/app"
)

func TestChainStart(t *testing.T) {
	a := app.New(log.NewNopLogger())
	if a == nil {
		t.Fatal("app should initialize")
	}
}
