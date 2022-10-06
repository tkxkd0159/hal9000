package base

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

type (
	AccAddr = sdktypes.AccAddress
)

func MakeTxFactory(ctx client.Context, gas string, gasPrice string, memo string, gasWeight float64) tx.Factory {
	gasSetting, _ := flags.ParseGasSetting(gas)

	initFac := tx.Factory{}.
		WithAccountNumber(0).
		WithSequence(0).
		WithTimeoutHeight(0).
		WithTxConfig(ctx.TxConfig).
		WithChainID(ctx.ChainID).
		WithKeybase(ctx.Keyring).
		WithAccountRetriever(ctx.AccountRetriever).
		WithSimulateAndExecute(gasSetting.Simulate).
		WithGas(gasSetting.Gas).
		WithGasPrices(gasPrice).
		WithGasAdjustment(flags.DefaultGasAdjustment * gasWeight)

	return initFac.
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT).
		WithMemo(memo)

}
