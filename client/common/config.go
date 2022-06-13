package common

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"io"
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
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

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

func MakeContext(mb module.BasicManager, from string, tmRPC string, chainID string, root string, backend string, userInput io.Reader) (client.Context, error) {
	encCfg := MakeEncodingConfig(mb)
	initClientCtx := client.Context{}.
		WithCodec(encCfg.Marshaler).
		WithInterfaceRegistry(encCfg.InterfaceRegistry).
		WithTxConfig(encCfg.TxConfig).
		WithLegacyAmino(encCfg.Amino).
		WithInput(userInput).
		WithSignModeStr(flags.SignModeDirect).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(root).
		WithKeyringDir(root).
		WithChainID(chainID).
		WithFrom(from).
		WithNodeURI(tmRPC)

	kb := MakeKeyring(initClientCtx, backend)
	initClientCtx = initClientCtx.WithKeyring(kb)

	tmClient, err := client.NewClientFromNode(tmRPC)
	if err != nil {
		return initClientCtx, err
	}
	return initClientCtx.WithClient(tmClient), nil

}
