package main

import (
	"flag"
	"fmt"
	"github.com/Carina-labs/HAL9000/client/base"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/tester"
	"github.com/Carina-labs/HAL9000/utils"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ibctypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"os"
	"sync"
)

var (
	wg      sync.WaitGroup
	ctx     client.Context
	botInfo keyring.Info
)

func init() {
	base.SetBechPrefix()
}

func main() {
	isTest := flag.Bool("test", false, "Decide whether it's test with localnet")
	keyname := flag.String("name", "nova_fake_bot", "Set unique key name (uid)")
	newacc := flag.Bool("add", false, "Start client with making new account")
	intv := flag.Int("interval", 10, "tx push interval (sec)")
	disp := flag.Bool("display", false, "Show context log through stdout")
	IBCChan := flag.String("chan", "channel-0", "Nova's channel")
	ZoneID := flag.String("host", "gaia", "hostchain name")
	flag.Parse()
	flags := cfg.FlagOpts{Test: *isTest, New: *newacc, Disp: *disp, Kn: *keyname, Period: *intv}

	// Open api endpoint to check bot
	wg.Add(2)

	cfg.SetChainInfo(flags.Test)
	Nova := &cfg.NovaInfo{}
	Nova.Set("fake_bot_addr", flags.Kn)
	krDir, logDir := cfg.SetInitialDir(flags.Kn, "logs/pusher")
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
		ctx = base.MakeContext(
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

		botInfo = base.MakeClientWithNewAcc(
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

		ctx = base.MakeContext(
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
		botInfo = base.LoadClientPubInfo(ctx, flags.Kn)
	}
	ctx = base.AddMoreFromInfo(ctx)
	txf := base.MakeTxFactory(ctx, "auto", "0unova", "", 1.1)

	// ###### Start target bot logic ######
	go func(interval int) {
		defer wg.Done()
		port := "transfer"
		DenomTrace := ibctypes.DenomTrace{BaseDenom: "uatom", Path: fmt.Sprintf("%s/%s", port, *IBCChan)}
		tester.DepositGal(ctx, txf, botInfo, interval, fdErr, tester.IBCInfo{ZoneID: *ZoneID, IBCChan: *IBCChan, IBCPort: port}, DenomTrace.IBCDenom(), 1000)
	}(flags.Period)

	wg.Wait()
}
