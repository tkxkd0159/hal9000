package cmd

import (
	"io"
	"os"

	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/Carina-labs/HAL9000/client/base"
	"github.com/Carina-labs/HAL9000/client/nova"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
)

func SetupBotBase(f cfg.BotCommon, krDir string, ctxOut io.Writer) (ctx client.Context, botInfo keyring.Info, txf tx.Factory) {
	flags := f.GetBase()
	base.SetBechPrefix()
	cfg.LoadChainInfo(flags.IsTest)
	NovaInfo := cfg.NewNovaInfo()
	BotScrt := cfg.NewBotScrt(NovaInfo.ChainID, "bot_addr", flags.Kn)

	if flags.New {
		nova.SetupBotKey(flags.Kn, krDir, NovaInfo, BotScrt)
		os.Exit(0)
	}

	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)
	os.Stdin = rpipe
	_, err = wpipe.Write([]byte(BotScrt.Passphrase()))
	utils.CheckErr(err, "", 0)

	ctx = base.MakeContext(
		novaapp.ModuleBasics,
		BotScrt.Address(),
		NovaInfo.TmRPC.String(),
		NovaInfo.ChainID,
		krDir,
		keyring.BackendFile,
		rpipe,
		ctxOut,
		false,
	)

	botInfo = base.LoadClientPubInfo(ctx, flags.Kn)
	ctx = base.AddMoreFromInfo(ctx)
	txf = base.MakeTxFactory(ctx, cfg.Gas, cfg.NovaGasPrice, "", cfg.GasWeight)
	return
}
