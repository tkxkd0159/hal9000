package main

import (
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/base/query"
	novatypes "github.com/Carina-labs/HAL9000/client/nova/types"
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
	flags := cfg.SetRestakeFlags()
	krDir, logDir := cfg.SetInitialDir(flags.Kn, flags.LogLocation)
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, cfg.StdLogFile, cfg.LocalErrlogFile, cfg.ExtRedirectErrlogFile, flags.Disp)
	defer utils.CloseFds(fdLog, fdErr, fdErrExt)
	ctx, krInfo, txf := cmd.SetupBotBase(flags, krDir, fdLog)

	wg := new(sync.WaitGroup)
	wg.Add(NumWorker)
	botch := make(chan time.Time)
	go api.OpenMonitoringSrv(wg, botch, flags)

	go func(interval int) {
		defer wg.Done()
		bot := novatypes.NewBot(ctx, txf, krInfo, flags.Period, fdErr, botch)
		hostZone := cfg.NewHostChainInfo(flags.HostChain)
		hostZone.Set()
		cq := query.NewCosmosQueryClient(hostZone.GrpcAddr)
		defer utils.CloseGrpc(cq.ClientConn)
		logic.IcaAutoStake(cq, bot, hostZone)
	}(flags.Period)

	wg.Wait()
}
