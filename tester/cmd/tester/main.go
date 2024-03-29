package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/tkxkd0159/HAL9000/api"
	"github.com/tkxkd0159/HAL9000/client"
	basetypes "github.com/tkxkd0159/HAL9000/client/base/types"
	novaq "github.com/tkxkd0159/HAL9000/client/nova/query"
	cfg "github.com/tkxkd0159/HAL9000/config"
	"github.com/tkxkd0159/HAL9000/logic"
	"github.com/tkxkd0159/HAL9000/utils"
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
	defer utils.CloseFds(fdLog, fdErr, fdErrExt)
	ctx, krInfo, txf, cni := client.SetupBotBase(flags, krDir, fdLog, cfg.ControlChain, "bot_addr")
	log.SetOutput(ctx.Output)

	wg := new(sync.WaitGroup)
	wg.Add(NumWorker)
	botch := make(chan time.Time)
	go api.OpenMonitoringSrv(wg, botch, flags)

	_ = basetypes.NewBot(botType, ctx, txf, krInfo, bf.Period, fdErr, botch)
	hzInfo := cfg.NewHostChainInfo(bf.HostChain)
	hzInfo.SetDetailInfos()
	hzInfo.WithIBCInfo(flags, botType)

	nq := novaq.NewNovaQueryClient(cni.GRPC.Host, cni.Secure)
	defer utils.CloseGrpc(nq.ClientConn)
	tmpseq := logic.FetchBotSeq(nq, cfg.ActWithdraw, "gaia")
	fmt.Println(tmpseq)
	/*
		cq := query.NewCosmosQueryClient(hostZone.GrpcAddr)
		defer utils.CloseGrpc(cq.ClientConn)
		delegatedToken, height, apphash := logic.OracleInfo(cq, hostZone.Validator)

		bot.Txf = bot.Txf.WithSequence(320)
		msg1 := novaTx.MakeMsgUpdateChainState(bot.KrInfo.GetAddress(), hostZone.Name, hostZone.Denom, delegatedToken, height, apphash)
		msgs := []sdktypes.Msg{msg1}
		_ = base.GenTxByBot(bot, false, msgs...)
		_ = cni
	*/
	wg.Wait()
}
