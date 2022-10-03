package main

import (
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/base"
	"github.com/Carina-labs/HAL9000/cmd"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/logic"
	"github.com/Carina-labs/HAL9000/utils"
	"sync"
	"time"
)

const (
	NumWorker = 2
)

func init() {
	base.SetBechPrefix()
}

func main() {

	flags := cfg.SetWithdrawFlags()
	krDir, logDir := cfg.SetInitialDir(flags.Kn, flags.LogLocation)
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, cfg.StdLogFile, cfg.LocalErrlogFile, cfg.ExtRedirectErrlogFile, flags.Disp)
	defer utils.CloseFds(fdLog, fdErr, fdErrExt)
	ctx, botInfo, txf := cmd.InitBaseBot(flags, krDir, fdLog)

	wg := new(sync.WaitGroup)
	wg.Add(NumWorker)
	botch := make(chan time.Time)
	go api.OpenMonitoringSrv(wg, botch, flags)

	// ###### Start target bot logic ######
	go func(interval int) {
		defer wg.Done()
		logic.UndelegateAndWithdraw(flags.HostChain, ctx, txf, botInfo, flags.HostIBC.Transfer, interval, fdErr, botch)
	}(flags.Period)

	wg.Wait()
}
