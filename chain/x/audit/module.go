package audit

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

// AppModuleBasic implements the basic methods for the audit module.
type AppModuleBasic struct{}

// Name returns the audit module name.
func (AppModuleBasic) Name() string { return "audit" }

// RegisterLegacyAminoCodec is a noop for this example.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// DefaultGenesis returns empty genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return json.RawMessage("{}")
}
