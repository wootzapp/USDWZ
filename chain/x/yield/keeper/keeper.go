package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

var yieldKey = []byte("yield")

// Accrue adds new yield amount.
func (k Keeper) Accrue(ctx sdk.Context, amt sdk.Int) {
	y := k.getInt(ctx)
	y = y.Add(amt)
	k.setInt(ctx, y)
}

// Distribute returns accrued yield and resets it.
func (k Keeper) Distribute(ctx sdk.Context) sdk.Int {
	y := k.getInt(ctx)
	k.setInt(ctx, sdk.ZeroInt())
	return y
}

// Accrued returns current accrued yield.
func (k Keeper) Accrued(ctx sdk.Context) sdk.Int { return k.getInt(ctx) }

func (k Keeper) getInt(ctx sdk.Context) sdk.Int {
	bz := ctx.KVStore(k.storeKey).Get(yieldKey)
	if bz == nil {
		return sdk.ZeroInt()
	}
	i, _ := sdk.NewIntFromString(string(bz))
	return i
}

func (k Keeper) setInt(ctx sdk.Context, v sdk.Int) {
	ctx.KVStore(k.storeKey).Set(yieldKey, []byte(v.String()))
}
