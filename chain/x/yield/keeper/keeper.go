package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/example/usdwz/chain/x/yield/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey *storetypes.KVStoreKey
}

func NewKeeper(cdc codec.BinaryCodec, key *storetypes.KVStoreKey) Keeper {
	return Keeper{cdc: cdc, storeKey: key}
}

func ModuleName() string { return types.ModuleName }
