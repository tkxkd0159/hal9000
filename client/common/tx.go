package common

import (
	"errors"
	"fmt"
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"os"
	"time"
)

type (
	AccAddr = sdktypes.AccAddress
)

func CheckAccAddr(target any) (AccAddr, error) {
	switch target := target.(type) {
	case AccAddr:
		return target, nil
	case string:
		addr, err := sdktypes.AccAddressFromBech32(target)
		if err != nil {
			return nil, err
		}
		return addr, nil
	case []byte:
		return target, nil
	default:
		return nil, errors.New("cannot covert target to AccAddress")
	}
}

// GenTxWithFactory
// 1. Generate a TX with Msg (TxBuilder). If you set --generate-only, it makes unsigned tx and never broadcast
// 2. Sign the generated transaction with the keyring's account
// 3. Broadcast the tx to the Tendermint node using gPRC
func GenTxWithFactory(errFd *os.File, ctx client.Context, txf tx.Factory, onlyGen bool, msgs ...sdktypes.Msg) {
	if onlyGen {
		ctx = ctx.WithGenerateOnly(true)
	}

	err := tx.GenerateOrBroadcastTxWithFactory(ctx, txf, msgs...)
	if err != nil {
		utils.LogErrWithFd(errFd, err, "something went wrong while make tx", 1)
	} else {
		_, err = ctx.Output.Write([]byte(fmt.Sprintf("%v: Tx was generated\n", time.Now())))
		utils.CheckErr(err, "cannot write log on output", 1)
	}
}

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
		WithGas(gasSetting.Gas).
		WithSimulateAndExecute(gasSetting.Simulate).
		WithGasAdjustment(flags.DefaultGasAdjustment * gasWeight).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	return initFac.
		WithGasPrices(gasPrice).
		WithMemo(memo)

}
