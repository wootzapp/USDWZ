package app

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// placeholder imports for future modules
	escrowkeeper "github.com/example/usdwz/chain/x/escrow/keeper"
	escrowtypes "github.com/example/usdwz/chain/x/escrow/types"
	stablecoinkeeper "github.com/example/usdwz/chain/x/stablecoin/keeper"
	stablecointypes "github.com/example/usdwz/chain/x/stablecoin/types"
	validatorkeeper "github.com/example/usdwz/chain/x/validator/keeper"
	validatortypes "github.com/example/usdwz/chain/x/validator/types"
	yieldkeeper "github.com/example/usdwz/chain/x/yield/keeper"
	yieldtypes "github.com/example/usdwz/chain/x/yield/types"
)

// UsdWzApp is the main application type.
type UsdWzApp struct {
	*baseapp.BaseApp
	codec codec.Codec

	// keys to access the multistore
	keys map[string]*storetypes.KVStoreKey

	// module keepers
	StablecoinKeeper stablecoinkeeper.Keeper
	EscrowKeeper     escrowkeeper.Keeper
	ValidatorKeeper  validatorkeeper.Keeper
	YieldKeeper      yieldkeeper.Keeper
}

// New creates a new UsdWzApp instance with basic modules wired.
func New(logger log.Logger) *UsdWzApp {
	cdc := codec.NewProtoCodec(types.NewInterfaceRegistry())
	bApp := baseapp.NewBaseApp("usdWz", logger, nil, nil)

	// create keys
	keys := map[string]*storetypes.KVStoreKey{
		stablecointypes.StoreKey: storetypes.NewKVStoreKey(stablecointypes.StoreKey),
		escrowtypes.StoreKey:     storetypes.NewKVStoreKey(escrowtypes.StoreKey),
		validatortypes.StoreKey:  storetypes.NewKVStoreKey(validatortypes.StoreKey),
		yieldtypes.StoreKey:      storetypes.NewKVStoreKey(yieldtypes.StoreKey),
	}

	app := &UsdWzApp{
		BaseApp: bApp,
		codec:   cdc,
		keys:    keys,
	}

	// initialize keepers
	app.StablecoinKeeper = stablecoinkeeper.NewKeeper(cdc, keys[stablecointypes.StoreKey])
	app.EscrowKeeper = escrowkeeper.NewKeeper(cdc, keys[escrowtypes.StoreKey])
	app.ValidatorKeeper = validatorkeeper.NewKeeper(cdc, keys[validatortypes.StoreKey])
	app.YieldKeeper = yieldkeeper.NewKeeper(cdc, keys[yieldtypes.StoreKey])

	return app
}

func (app *UsdWzApp) ModuleManager() *module.Manager {
	return module.NewManager()
}

// Name returns the application name.
func (app *UsdWzApp) Name() string { return "usdWz" }

// Codec returns the app codec.
func (app *UsdWzApp) Codec() codec.Codec { return app.codec }

// KVStoreKey allows modules to access their stores.
func (app *UsdWzApp) KVStoreKey(name string) *storetypes.KVStoreKey { return app.keys[name] }
