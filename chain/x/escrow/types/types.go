package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName = "escrow"
	StoreKey   = ModuleName
)

// Escrow represents a simple escrow record.
type Escrow struct {
	Amount    sdk.Int `json:"amount"`
	Completed bool    `json:"completed"`
}
