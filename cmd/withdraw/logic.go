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

func UndelegateAndWithdraw(host string, txf tx.Factory, chanID string, interval int, errLogger *os.File) {

	targetIP := viper.GetString(fmt.Sprintf("net.ip.%s", host))
	targetGrpcAddr := targetIP + ":" + viper.GetString("net.port.grpc")
	conn, err := grpc.Dial(
		targetGrpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", 0)
	defer func(c *grpc.ClientConn) {
		err = c.Close()
		utils.CheckErr(err, "", 1)
	}(conn)
	cq := &query.CosmosQueryClient{ClientConn: conn}

	isStart := true
	stream := ut.Fstream{Err: errLogger}
	i := 0
	intv := time.Duration(interval)
	for {
		log.Printf("Undelegate & Withdraw Bot is ongoing for %d secs\n", int(intv)*i)

		secs := cq.GetLatestBlock().GetBlock().GetHeader().GetTime().GetSeconds()
		nanos := cq.GetLatestBlock().GetBlock().GetHeader().GetTime().GetNanos()
		currentTs := time.Unix(secs, int64(nanos)).UTC()

		if isStart {
			msg1 := novaTx.MakeMsgUndelegate(host, botInfo.GetAddress())
			msgs := []sdktypes.Msg{msg1}
			common.GenTxWithFactory(stream, ctx, txf, false, msgs...)
			isStart = false
		} else {
			msg1 := novaTx.MakeMsgUndelegate(host, botInfo.GetAddress())
			msgs := []sdktypes.Msg{msg1}
			common.GenTxWithFactory(stream, ctx, txf, false, msgs...)
			time.Sleep(60 * time.Second)

			msg2 := novaTx.MakeMsgPendingWithdraw(host, botInfo.GetAddress(), "transfer", chanID, currentTs)
			msgs = []sdktypes.Msg{msg2}
			common.GenTxWithFactory(stream, ctx, txf, false, msgs...)
		}

		time.Sleep(intv * time.Second)
		i++
	}

}
