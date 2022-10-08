package main

import (
	"sync"
	"time"

	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/base/query"
	novatypes "github.com/Carina-labs/HAL9000/client/nova/types"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/logic"
	"github.com/Carina-labs/HAL9000/utils"
)

const (
	NumWorker = 2
)

func main() {
	flags := cfg.SetFlags(cfg.ActRestake)
	bf := flags.GetBase()

	krDir, logDir := cfg.SetInitialDir(bf.Kn, bf.LogLocation)
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, cfg.StdLogFile, cfg.LocalErrlogFile, cfg.ExtRedirectErrlogFile, bf.Disp)
	defer utils.CloseFds(fdLog, fdErr, fdErrExt)
	ctx, krInfo, txf := cfg.SetupBotBase(flags, krDir, fdLog)

	wg := new(sync.WaitGroup)
	wg.Add(NumWorker)
	botch := make(chan time.Time)
	go api.OpenMonitoringSrv(wg, botch, flags)

	go func() {
		defer wg.Done()
		bot := novatypes.NewBot(ctx, txf, krInfo, bf.Period, fdErr, botch)
		hostZone := cfg.NewHostChainInfo(bf.HostChain)
		hostZone.Set()
		cq := query.NewCosmosQueryClient(hostZone.GrpcAddr)
		defer utils.CloseGrpc(cq.ClientConn)
		logic.IcaAutoStake(cq, bot, hostZone)
	}()

	wg.Wait()
}
