package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/example/usdwz/chain/x/validator/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey *storetypes.KVStoreKey
}

func NewKeeper(cdc codec.BinaryCodec, key *storetypes.KVStoreKey) Keeper {
	return Keeper{cdc: cdc, storeKey: key}
}

func ModuleName() string { return types.ModuleName }

// SubmitVote records a validator vote.
func (k Keeper) SubmitVote(ctx sdk.Context, val string, approve bool) {
	ctx.KVStore(k.storeKey).Set([]byte("vote:"+val), []byte(strconv.FormatBool(approve)))
}

// Vote returns a stored vote.
func (k Keeper) Vote(ctx sdk.Context, val string) (bool, bool) {
	bz := ctx.KVStore(k.storeKey).Get([]byte("vote:" + val))
	if bz == nil {
		return false, false
	}
	b, err := strconv.ParseBool(string(bz))
	if err != nil {
		return false, false
	}
	return b, true
}

// Tally returns counts of yes and no votes.
func (k Keeper) Tally(ctx sdk.Context) (yes, no int) {
	it := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte("vote:"))
	defer it.Close()
	for ; it.Valid(); it.Next() {
		b, _ := strconv.ParseBool(string(it.Value()))
		if b {
			yes++
		} else {
			no++
		}
	}
	return
}
