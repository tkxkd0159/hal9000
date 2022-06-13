package common

import (
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"io"
	"os"
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

func MakeContext(mb module.BasicManager, from string, tmRPC string, chainID string, root string, backend string, userInput io.Reader) (client.Context, error) {
	encCfg := MakeEncodingConfig(mb)
	initClientCtx := client.Context{}.
		WithSimulation(false).
		WithInput(userInput).
		WithCodec(encCfg.Marshaler).
		WithInterfaceRegistry(encCfg.InterfaceRegistry).
		WithTxConfig(encCfg.TxConfig).
		WithLegacyAmino(encCfg.Amino).
		WithSignModeStr(flags.SignModeDirect).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(root).
		WithKeyringDir(root).
		WithChainID(chainID).
		WithNodeURI(tmRPC).
		WithFrom(from).
		WithOutput(os.Stdout)

	kb := MakeKeyring(initClientCtx, backend)
	initClientCtx = initClientCtx.WithKeyring(kb)

	tmClient, err := client.NewClientFromNode(tmRPC)
	utils.CheckErr(err, "-> Cannot set node client", 0)
	if err != nil {
		return initClientCtx, err
	}
	return initClientCtx.
		WithClient(tmClient).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithSkipConfirmation(true), nil
}

func AddMoreFromInfo(ctx client.Context) client.Context {
	fromAddr, fromName, _, err := client.GetFromFields(ctx.Keyring, ctx.From, ctx.GenerateOnly)
	utils.CheckErr(err, "cannot get info from keyring", 0)
	ctx = ctx.WithFromAddress(fromAddr).WithFromName(fromName)
	return ctx
}

// gas = "auto", fee = "0unova", gasPrice = "ounova"
func MakeTxFactory(ctx client.Context, gas string, gasPrice string, memo string) tx.Factory {
	gasSetting, _ := flags.ParseGasSetting(gas)

	initFac := tx.Factory{}.
		WithAccountNumber(0).
		WithSequence(0).
		WithTimeoutHeight(0).
		WithTxConfig(ctx.TxConfig).
		WithChainID(ctx.ChainID).
		WithKeybase(ctx.Keyring).
		WithAccountRetriever(ctx.AccountRetriever).
		WithGas(gasSetting.Gas).
		WithSimulateAndExecute(gasSetting.Simulate).
		WithGasAdjustment(flags.DefaultGasAdjustment).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	return initFac.
		WithGasPrices(gasPrice).
		WithMemo(memo)

}
