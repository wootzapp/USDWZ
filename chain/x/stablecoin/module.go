package stablecoin

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"

	"github.com/example/usdwz/chain/x/stablecoin/keeper"
	"github.com/example/usdwz/chain/x/stablecoin/types"
)

// AppModuleBasic implements basic methods for the module.
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string                                   { return types.ModuleName }
func (AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino)    {}
func (AppModuleBasic) DefaultGenesis(codec.JSONCodec) json.RawMessage { return json.RawMessage("{}") }
func (AppModuleBasic) ValidateGenesis(codec.JSONCodec, client.TxEncodingConfig, json.RawMessage) error {
	return nil
}
func (AppModuleBasic) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {}
func (AppModuleBasic) GetTxCmd() *cobra.Command                                    { return nil }
func (AppModuleBasic) GetQueryCmd() *cobra.Command                                 { return nil }

// AppModule wraps the keeper.
type AppModule struct {
	AppModuleBasic
	Keeper keeper.Keeper
}

// NewAppModule creates a new AppModule.
func NewAppModule(k keeper.Keeper) AppModule { return AppModule{Keeper: k} }
