package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/example/usdwz/chain/x/stablecoin/types"
)

// Keeper handles all stablecoin state changes.
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey *storetypes.KVStoreKey
}

// NewKeeper returns a new Keeper.
func NewKeeper(cdc codec.BinaryCodec, key *storetypes.KVStoreKey) Keeper {
	return Keeper{cdc: cdc, storeKey: key}
}

// ModuleName returns the module name.
func ModuleName() string { return types.ModuleName }
