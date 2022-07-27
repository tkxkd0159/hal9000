package logic

import (
	"fmt"
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/client/common/query"
	novaTx "github.com/Carina-labs/HAL9000/client/nova/msgs"
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func UpdateChainState(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File) {
	var (
		targetIP       = viper.GetString(fmt.Sprintf("net.ip.%s", host))
		targetGrpcAddr = targetIP + ":" + viper.GetString("net.port.grpc")
		targetValAddr  = viper.GetString(fmt.Sprintf("%s.val_addr", host))
		targetDenom    = viper.GetString(fmt.Sprintf("%s.denom", host))
		targetDecimal  = viper.GetUint32(fmt.Sprintf("%s.decimal", host))
	)

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

		msg1 := novaTx.MakeMsgUpdateChainState(botInfo.GetAddress(), host, delegatedToken, targetDenom, targetDecimal, h, apphash)
		msgs := []sdktypes.Msg{msg1}
		common.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaAutoStake(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File) {
	var (
		targetIP       = viper.GetString(fmt.Sprintf("net.ip.%s", host))
		targetGrpcAddr = targetIP + ":" + viper.GetString("net.port.grpc")
		targetValAddr  = viper.GetString(fmt.Sprintf("%s.val_addr", host))
		targetHostAddr = viper.GetString(fmt.Sprintf("%s.host_addr", host))
	)

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

	i := 0
	intv := time.Duration(interval)
	for {
		log.Printf("Re-staking Bot is ongoing for %d secs\n", int(intv)*i)

		r := cq.GetRewards(targetHostAddr, targetValAddr).GetRewards()[0]
		msg1 := novaTx.MakeMsgIcaAutoStaking(host, targetHostAddr, botInfo.GetAddress(), r)
		msgs := []sdktypes.Msg{msg1}
		common.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaStake(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, chanID string, interval int, errLogger *os.File) {

	i := 0
	intv := time.Duration(interval)
	for {
		log.Printf("ICA-staking Bot is ongoing for %d secs\n", int(intv)*i)

		msg1 := novaTx.MakeMsgDelegate(host, botInfo.GetAddress(), "transfer", chanID)
		msgs := []sdktypes.Msg{msg1}
		common.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}

func UndelegateAndWithdraw(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, chanID string, interval int, errLogger *os.File) {

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
			common.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
			isStart = false
		} else {
			msg1 := novaTx.MakeMsgUndelegate(host, botInfo.GetAddress())
			msgs := []sdktypes.Msg{msg1}
			common.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)

			time.Sleep(60 * time.Second)

			msg2 := novaTx.MakeMsgPendingWithdraw(host, botInfo.GetAddress(), "transfer", chanID, currentTs)
			msgs = []sdktypes.Msg{msg2}
			common.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		}

		time.Sleep(intv * time.Second)
		i++
	}

}
