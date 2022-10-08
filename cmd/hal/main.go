package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Carina-labs/HAL9000/api"
	novatypes "github.com/Carina-labs/HAL9000/client/nova/types"
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
	fmt.Printf("%#v", flags)

	bf := flags.GetBase()
	fmt.Printf("%#v\n", flags)
	krDir, logDir := cfg.SetInitialDir(bf.Kn, bf.LogLocation)
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, cfg.StdLogFile, cfg.LocalErrlogFile, cfg.ExtRedirectErrlogFile, bf.Disp)
	defer utils.CloseFds(fdLog, fdErr, fdErrExt)
	ctx, krInfo, txf := cfg.SetupBotBase(flags, krDir, fdLog)

	wg := new(sync.WaitGroup)
	wg.Add(NumWorker)
	botch := make(chan time.Time)
	go api.OpenMonitoringSrv(wg, botch, flags)

	bot := novatypes.NewBot(ctx, txf, krInfo, bf.Period, fdErr, botch)
	hostZone := cfg.NewHostChainInfo(bf.HostChain)
	hostZone.Set()
	hostZone.WithIBCInfo(flags, botType)
	go func() {
		defer wg.Done()
		logic.RouteBotAction(botType, bot, hostZone)
	}()

	wg.Wait()
}
