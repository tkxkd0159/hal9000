package main

import (
	"github.com/Carina-labs/HAL9000/api"
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

func main() {
	flags := cfg.SetOracleFlags()
	krDir, logDir := cfg.SetInitialDir(flags.Kn, flags.LogLocation)
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, cfg.StdLogFile, cfg.LocalErrlogFile, cfg.ExtRedirectErrlogFile, flags.Disp)
	defer utils.CloseFds(fdLog, fdErr, fdErrExt)
	ctx, botInfo, txf := cmd.InitBaseBot(flags, krDir, fdLog)

	wg := new(sync.WaitGroup)
	wg.Add(NumWorker)
	botch := make(chan time.Time)
	go api.OpenMonitoringSrv(wg, botch, flags)

	go func(interval int) {
		defer wg.Done()
		logic.UpdateChainState(flags.HostChain, ctx, txf, botInfo, interval, fdErr, botch)
	}(flags.Period)

	wg.Wait()
}
