package keeper

import (
	"errors"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/example/usdwz/chain/monitoring"
	"github.com/example/usdwz/chain/x/stablecoin/types"
)

// Keeper handles all stablecoin state changes.
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey *storetypes.KVStoreKey
}

var (
	supplyKey     = []byte("supply")
	collateralKey = []byte("collateral")
)

// NewKeeper returns a new Keeper.
func NewKeeper(cdc codec.BinaryCodec, key *storetypes.KVStoreKey) Keeper {
	return Keeper{cdc: cdc, storeKey: key}
}

// ModuleName returns the module name.
func ModuleName() string { return types.ModuleName }

// DepositCollateral increases collateral backing and updates metrics.
func (k Keeper) DepositCollateral(ctx sdk.Context, amt sdk.Int) {
	coll := k.getInt(ctx, collateralKey)
	coll = coll.Add(amt)
	k.setInt(ctx, collateralKey, coll)
	monitoring.SetCollateral(float64(coll.Int64()))
}

// RedeemCollateral decreases collateral if sufficient backing remains.
func (k Keeper) RedeemCollateral(ctx sdk.Context, amt sdk.Int) error {
	coll := k.getInt(ctx, collateralKey)
	supply := k.getInt(ctx, supplyKey)
	if coll.Sub(amt).LT(supply) {
		return errors.New("insufficient collateral")
	}
	coll = coll.Sub(amt)
	k.setInt(ctx, collateralKey, coll)
	monitoring.SetCollateral(float64(coll.Int64()))
	return nil
}

// Mint creates new USDWZ if backed by collateral.
func (k Keeper) Mint(ctx sdk.Context, amt sdk.Int) error {
	supply := k.getInt(ctx, supplyKey)
	coll := k.getInt(ctx, collateralKey)
	if supply.Add(amt).GT(coll) {
		return errors.New("insufficient collateral")
	}
	supply = supply.Add(amt)
	k.setInt(ctx, supplyKey, supply)
	return nil
}

// Burn removes USDWZ from supply.
func (k Keeper) Burn(ctx sdk.Context, amt sdk.Int) error {
	supply := k.getInt(ctx, supplyKey)
	if supply.LT(amt) {
		return errors.New("insufficient supply")
	}
	supply = supply.Sub(amt)
	k.setInt(ctx, supplyKey, supply)
	return nil
}

// GetSupply returns current USDWZ supply.
func (k Keeper) GetSupply(ctx sdk.Context) sdk.Int { return k.getInt(ctx, supplyKey) }

// GetCollateral returns current collateral amount.
func (k Keeper) GetCollateral(ctx sdk.Context) sdk.Int {
	return k.getInt(ctx, collateralKey)
}

// TestContext creates a context for keeper tests.
func TestContext(key *storetypes.KVStoreKey) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()
	return sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())
}

func (k Keeper) getInt(ctx sdk.Context, key []byte) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroInt()
	}
	i, _ := sdk.NewIntFromString(string(bz))
	return i
}

func (k Keeper) setInt(ctx sdk.Context, key []byte, v sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(key, []byte(v.String()))
}
