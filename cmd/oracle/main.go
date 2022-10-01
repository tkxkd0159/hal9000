package main

import (
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/base"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/logic"
	"github.com/Carina-labs/HAL9000/utils"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"os"
	"sync"
	"time"
)

var (
	ctx     client.Context
	botInfo keyring.Info
)

const (
	NumWorker = 2
)

func init() {
	base.SetBechPrefix()
}

func main() {
	flags := cfg.SetOracleFlags()
	cfg.LoadChainInfo(flags.IsTest)
	NovaInfo := cfg.NewNovaInfo().Set("bot_addr", flags.Kn)
	krDir, logDir := cfg.SetInitialDir(flags.Kn, flags.LogLocation)
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, cfg.StdLogFile, cfg.LocalErrlogFile, cfg.ExtRedirectErrlogFile, flags.Disp)
	defer utils.CloseFds(fdLog, fdErr, fdErrExt)

	if flags.New {
		cfg.SetupBotKey(flags.Kn, krDir, NovaInfo)
		os.Exit(0)
	}
	wg := new(sync.WaitGroup)
	wg.Add(NumWorker)
	botch := make(chan time.Time)
	go api.OpenMonitoringSrv(wg, botch, flags)

	// set pipe to ignore stdin tty
	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)
	os.Stdin = rpipe
	_, err = wpipe.Write([]byte(NovaInfo.Bot.Passphrase()))
	utils.CheckErr(err, "", 0)

	ctx = base.MakeContext(
		novaapp.ModuleBasics,
		NovaInfo.Bot.Addr,
		NovaInfo.TmRPC.String(),
		NovaInfo.ChainID,
		krDir,
		keyring.BackendFile,
		rpipe,
		fdLog,
		false,
	)

	botInfo = base.LoadClientPubInfo(ctx, flags.Kn)
	ctx = base.AddMoreFromInfo(ctx)
	txf := base.MakeTxFactory(ctx, cfg.Gas, cfg.NovaGasPrice, "", cfg.GasWeight)

	// ###### Start target bot logic ######
	go func(interval int) {
		defer wg.Done()
		logic.UpdateChainState(flags.HostChain, ctx, txf, botInfo, interval, fdErr, botch)
	}(flags.Period)

	wg.Wait()
}
