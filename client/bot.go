package client

import (
	"fmt"
	"io"
	"log"
	"os"

	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"

	"github.com/Carina-labs/HAL9000/client/base"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
)

func SetupBotBase(f cfg.BotCommon, krDir string, ctxOut io.Writer, zone string, target string) (ctx client.Context, botInfo keyring.Info, txf tx.Factory, cni *cfg.ChainNetInfo) {
	flags := f.GetBase()
	base.SetBechPrefix()
	cfg.LoadChainInfo(flags.IsTest)
	cni = cfg.NewChainNetInfo(zone)
	BotScrt := NewBotScrt(cni.ChainID, target, flags.Kn)

	if flags.New {
		SetupBotKey(flags.Kn, krDir, cni, BotScrt)
		log.Println("🎉 Your keyring has been successfully set.")
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
		cni.TmRPC.String(),
		cni.ChainID,
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

func SetupBotKey(keyname, keyloc string, info *cfg.ChainNetInfo, bot BotScrt) {
	ctx := base.MakeContext(
		novaapp.ModuleBasics,
		bot.Address(),
		info.TmRPC.String(),
		info.ChainID,
		keyloc,
		keyring.BackendFile,
		os.Stdin,
		os.Stdout,
		false,
	)

	_ = base.MakeClientWithNewAcc(
		ctx,
		keyname,
		cfg.InputMnemonic(),
		sdktypes.FullFundraiserPath,
		hd.Secp256k1,
	)
}

type BotScrt struct {
	addr       string
	passphrase string
}

func NewBotScrt(zone string, addrTarget string, keyname ...string) (bi BotScrt) {
	if len(keyname) == 1 {
		bi.passphrase = GetPassphrase(cfg.Sviper)
	}
	bi.addr = viper.GetString(fmt.Sprintf("%s.%s", zone, addrTarget))
	return
}

func (b BotScrt) Address() string {
	return b.addr
}

func (b BotScrt) Passphrase() string {
	return b.passphrase
}

func GetPassphrase(vp *viper.Viper) string {
	pw := vp.GetString("pw")
	pp := fmt.Sprintf("%s\n%s\n", pw, pw)
	return pp
}
