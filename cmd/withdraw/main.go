package main

import (
	"flag"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/logic"
	"github.com/Carina-labs/HAL9000/utils"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"os"
	"sync"
)

var (
	wg      sync.WaitGroup
	ctx     client.Context
	botInfo keyring.Info
)

func init() {
	common.SetBechPrefix()
}

func main() {
	isTest := flag.Bool("test", false, "Decide whether it's test with localnet")
	apiAddr := flag.String("api", "127.0.0.1:3337", "Set bot api address")
	keyname := flag.String("name", "nova_bot", "Set unique key name (uid)")
	newacc := flag.Bool("add", false, "Start client with making new account")
	chanID := flag.String("ch", "channel-45", "Host Transfer Channel ID")
	hostchain := flag.String("host", "gaia", "Name of the host chain from which to obtain oracle info")
	intv := flag.Int("interval", 21*24*60*60, "Withdraw interval (sec)")
	disp := flag.Bool("display", false, "Show context log through stdout")
	flag.Parse()
	flags := cfg.FlagOpts{Test: *isTest, New: *newacc, Disp: *disp, ExtIP: *apiAddr, Kn: *keyname, Host: *hostchain, Period: *intv, IBCChan: cfg.IBCChan{Host: cfg.IBCPort{Transfer: *chanID}}}

	wg.Add(2)
	go func() {
		defer wg.Done()
		api.Server{}.On(flags.ExtIP)
	}()

	cfg.SetChainInfo(flags.Test)
	Nova := &cfg.NovaInfo{}
	Nova.Set("bot_addr", flags.Kn)
	krDir, logDir := cfg.SetInitialDir(flags.Kn, "logs/withdraw")
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, "ctxlog.txt", "nova_err.txt", "other_err.txt", flags.Disp)
	projFps := []*os.File{fdLog, fdErr, fdErrExt}
	defer func(fps ...*os.File) {
		for _, fp := range fps {
			err := fp.Close()
			utils.CheckErr(err, "", 1)
		}
	}(projFps...)

	// set pipe to ignore stdin tty
	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)

	if flags.New {
		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			Nova.Bot.Addr,
			Nova.TmRPC.String(),
			Nova.ChainID,
			krDir,
			keyring.BackendFile,
			os.Stdin,
			fdLog,
			false,
		)
		botInfo = common.MakeClientWithNewAcc(
			ctx,
			flags.Kn,
			Nova.Bot.Mnemonic(),
			sdktypes.FullFundraiserPath,
			hd.Secp256k1,
		)
		os.Exit(0)
	} else {
		pp := Nova.Bot.Passphrase()
		_, err = wpipe.Write([]byte(pp))
		utils.CheckErr(err, "", 0)

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			Nova.Bot.Addr,
			Nova.TmRPC.String(),
			Nova.ChainID,
			krDir,
			keyring.BackendFile,
			rpipe,
			fdLog,
			false,
		)
		os.Stdin = rpipe
		botInfo = common.LoadClientPubInfo(ctx, flags.Kn)
	}
	ctx = common.AddMoreFromInfo(ctx)
	txf := common.MakeTxFactory(ctx, "auto", "0unova", "", 1.1)

	// ###### Start target bot logic ######
	go func(interval int) {
		defer wg.Done()
		logic.UndelegateAndWithdraw(flags.Host, ctx, txf, botInfo, flags.IBCChan.Host.Transfer, interval, fdErr)
	}(flags.Period)

	wg.Wait()
}
