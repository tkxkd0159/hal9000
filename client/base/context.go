package base

import (
	"io"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/Carina-labs/HAL9000/utils"
)

func MakeContext(mb module.BasicManager, from string, tmRPC string, chainID string, root string, backend string, userInput io.Reader, userOutput io.Writer, genOnly bool) client.Context {
	encCfg := MakeEncodingConfig(mb)
	initClientCtx := client.Context{}.
		WithSimulation(false).
		WithSkipConfirmation(true).
		WithSignModeStr(flags.SignModeDirect).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithFrom(from).
		WithNodeURI(tmRPC).
		WithChainID(chainID).
		WithHomeDir(root).
		WithKeyringDir(root).
		WithInput(userInput).
		WithOutput(userOutput).
		WithGenerateOnly(genOnly).
		WithCodec(encCfg.Marshaler).
		WithInterfaceRegistry(encCfg.InterfaceRegistry).
		WithTxConfig(encCfg.TxConfig).
		WithLegacyAmino(encCfg.Amino)

	kb := MakeKeyring(initClientCtx, backend)
	tmClient, err := client.NewClientFromNode(tmRPC)
	utils.CheckErr(err, "-> Cannot set node client", 0)

	return initClientCtx.
		WithKeyring(kb).
		WithClient(tmClient)

}

func AddMoreFromInfo(ctx client.Context) client.Context {
	fromAddr, fromName, _, err := client.GetFromFields(ctx, ctx.Keyring, ctx.From)
	utils.CheckErr(err, "cannot get info from keyring", 0)
	ctx = ctx.WithFromAddress(fromAddr).WithFromName(fromName)
	return ctx
}
