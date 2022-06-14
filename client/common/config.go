package common

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

func makeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	ir := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(ir)
	txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: ir,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}

func MakeEncodingConfig(mb module.BasicManager) EncodingConfig {
	encCfg := makeEncodingConfig()
	std.RegisterLegacyAminoCodec(encCfg.Amino)
	std.RegisterInterfaces(encCfg.InterfaceRegistry)
	mb.RegisterLegacyAminoCodec(encCfg.Amino)
	mb.RegisterInterfaces(encCfg.InterfaceRegistry)

	return encCfg
}
