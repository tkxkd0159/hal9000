package base

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	txcore "github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
)

type (
	AccAddr = sdktypes.AccAddress
)

func MakeTxFactory(ctx client.Context, gas string, gasPrice string, memo string, gasWeight float64) txcore.Factory {
	gasSetting, _ := flags.ParseGasSetting(gas)

	initFac := txcore.Factory{}.
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

// BroadcastTx broadcast to a Tendermint node
func BroadcastTx(ctx client.Context, txBytes []byte) error {
	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}
	return ctx.PrintProto(res)
}

// GenerateTx make txbytes
func GenerateTx(ctx client.Context, txf txcore.Factory, msgs ...sdktypes.Msg) ([]byte, error) {
	txf, err := prepareFactory(ctx, txf)
	if err != nil {
		return nil, err
	}

	if txf.SimulateAndExecute() || ctx.Simulate {
		_, adjusted, err := txcore.CalculateGas(ctx, txf, msgs...)
		if err != nil {
			return nil, err
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", txcore.GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	if ctx.Simulate {
		return nil, nil
	}

	tx, err := txcore.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}

	if !ctx.SkipConfirm {
		out, err := ctx.TxConfig.TxJSONEncoder()(tx.GetTx())
		if err != nil {
			return nil, err
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return nil, err
		}
	}

	tx.SetFeeGranter(ctx.GetFeeGranterAddress())
	err = txcore.Sign(txf, ctx.GetFromName(), tx, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := ctx.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return nil, err
	}
	return txBytes, nil
}

func prepareFactory(clientCtx client.Context, txf txcore.Factory) (txcore.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}
