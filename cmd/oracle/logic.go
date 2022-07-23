package main

import (
	"fmt"
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/client/common/query"
	novaTx "github.com/Carina-labs/HAL9000/client/nova/msgs"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func UpdateChainState(host *string, txf tx.Factory, interval int, errLogger *os.File) {
	var (
		targetIP       = viper.GetString(fmt.Sprintf("net.ip.%s", *host))
		targetGrpcAddr = targetIP + ":" + viper.GetString("net.port.grpc")
		targetValAddr  = viper.GetString(fmt.Sprintf("%s.val_addr", *host))
		targetDenom    = viper.GetString(fmt.Sprintf("%s.denom", *host))
		targetDecimal  = viper.GetUint32(fmt.Sprintf("%s.decimal", *host))
	)

	conn, err := grpc.Dial(
		targetGrpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", 0)
	defer func(c *grpc.ClientConn) {
		err = c.Close()
		if err != nil {
			log.Printf("unexpected gRPC disconnection: %v", err)
		}
	}(conn)
	cq := &query.CosmosQueryClient{ClientConn: conn}

	stream := ut.Fstream{Err: errLogger}
	i := 0
	intv := time.Duration(interval)
	for {
		log.Printf("Bot is ongoing for %d secs\n", int(intv)*i)

		h := cq.GetLatestBlock().GetBlock().GetHeader().GetHeight()
		hisInfo := cq.GetHistoricalInfo(h)
		apphash := hisInfo.GetHist().GetHeader().GetAppHash()

		var delegatedToken string
		for _, val := range hisInfo.GetHist().GetValset() {
			if val.GetOperatorAddress() != targetValAddr {
				continue
			} else {
				delegatedToken = val.GetTokens()
				break
			}
		}

		msg1 := novaTx.MakeMsgUpdateChainState(botInfo.GetAddress(), *host, delegatedToken, targetDenom, targetDecimal, h, apphash)
		msgs := []sdktypes.Msg{msg1}
		common.GenTxWithFactory(stream, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}
