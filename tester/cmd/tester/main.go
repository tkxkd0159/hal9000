package main

import (
	"sync"
	"time"

	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/base/query"
	basetypes "github.com/Carina-labs/HAL9000/client/base/types"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/logic"
	"github.com/Carina-labs/HAL9000/utils"
)

const (
	NumWorker = 2
)

// TODO: add subcommand for testing
func main() {
	flags := cfg.SetFlags(cfg.ActOracle)
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
		bot := basetypes.NewBot(ctx, txf, krInfo, bf.Period, fdErr, botch)
		hostZone := cfg.NewHostChainInfo(bf.HostChain)
		hostZone.Set()
		cq := query.NewCosmosQueryClient(hostZone.GrpcAddr)
		defer utils.CloseGrpc(cq.ClientConn)
		logic.UpdateChainState(cq, bot, hostZone)
	}()

	wg.Wait()
}
