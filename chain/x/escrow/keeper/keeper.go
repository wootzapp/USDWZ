package keeper

import (
	"encoding/json"
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/example/usdwz/chain/x/escrow/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey *storetypes.KVStoreKey
}

func NewKeeper(cdc codec.BinaryCodec, key *storetypes.KVStoreKey) Keeper {
	return Keeper{cdc: cdc, storeKey: key}
}

func ModuleName() string { return types.ModuleName }

// CreateEscrow stores a new escrow entry.
func (k Keeper) CreateEscrow(ctx sdk.Context, id string, amt sdk.Int) {
	e := types.Escrow{Amount: amt, Completed: false}
	bz, _ := json.Marshal(e)
	ctx.KVStore(k.storeKey).Set([]byte("escrow:"+id), bz)
}

// FinalizeEscrow marks an escrow as completed.
func (k Keeper) FinalizeEscrow(ctx sdk.Context, id string) error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("escrow:" + id))
	if bz == nil {
		return errors.New("not found")
	}
	var e types.Escrow
	_ = json.Unmarshal(bz, &e)
	e.Completed = true
	bz, _ = json.Marshal(e)
	store.Set([]byte("escrow:"+id), bz)
	return nil
}

// GetEscrow retrieves an escrow record by id.
func (k Keeper) GetEscrow(ctx sdk.Context, id string) (types.Escrow, bool) {
	bz := ctx.KVStore(k.storeKey).Get([]byte("escrow:" + id))
	if bz == nil {
		return types.Escrow{}, false
	}
	var e types.Escrow
	_ = json.Unmarshal(bz, &e)
	return e, true
}

// DeleteEscrow removes an escrow record.
func (k Keeper) DeleteEscrow(ctx sdk.Context, id string) {
	ctx.KVStore(k.storeKey).Delete([]byte("escrow:" + id))
}

// ListEscrows returns all escrows.
func (k Keeper) ListEscrows(ctx sdk.Context) []types.Escrow {
	store := ctx.KVStore(k.storeKey)
	it := sdk.KVStorePrefixIterator(store, []byte("escrow:"))
	defer it.Close()
	var out []types.Escrow
	for ; it.Valid(); it.Next() {
		var e types.Escrow
		if err := json.Unmarshal(it.Value(), &e); err == nil {
			out = append(out, e)
		}
	}
	return out
}
