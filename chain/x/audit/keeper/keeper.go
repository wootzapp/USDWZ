package keeper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/example/usdwz/chain/x/audit/types"
)

// Keeper manages audit attestations.
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

// PublishAttestation stores an attestation under the current block height.
func (k Keeper) PublishAttestation(ctx sdk.Context, hash string) error {
	att := types.Attestation{Timestamp: ctx.BlockTime().Unix(), Hash: hash}
	bz, err := json.Marshal(att)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	key := []byte("attestation:" + hash)
	store.Set(key, bz)
	return nil
}
