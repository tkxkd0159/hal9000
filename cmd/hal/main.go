package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/Carina-labs/HAL9000/api"
	basetypes "github.com/Carina-labs/HAL9000/client/base/types"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/logic"
	"github.com/Carina-labs/HAL9000/utils"
)

const (
	NumWorker = 2
)

func main() {
	botType := cfg.CheckBotType(os.Args[1])
	flags := cfg.SetFlags(botType)
	bf := flags.GetBase()
	krDir, logDir := cfg.SetInitialDir(bf.Kn, bf.LogLocation)
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, cfg.StdLogFile, cfg.LocalErrlogFile, cfg.ExtRedirectErrlogFile, bf.Disp)
	if !bf.Disp {
		defer utils.CloseFds(fdLog, fdErr, fdErrExt)
	}

	botctx, krInfo, txf, cni := cfg.SetupBotBase(flags, krDir, fdLog, cfg.ControlChain, cfg.NovaBotAddrKey)
	log.SetOutput(botctx.Output)

	wg := new(sync.WaitGroup)
	wg.Add(NumWorker)
	botch := make(chan time.Time)
	go api.OpenMonitoringSrv(wg, botch, flags)
	bot := basetypes.NewBot(botctx, txf, krInfo, bf.Period, fdErr, botch)
	hostZone := cfg.NewHostChainInfo(bf.HostChain)
	hostZone.Set()
	hostZone.WithIBCInfo(flags, botType)
	go func() {
		defer wg.Done()
		logic.RouteBotAction(botType, bot, cni, hostZone)
	}()

	wg.Wait()
}
