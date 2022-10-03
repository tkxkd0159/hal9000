package logic

import (
	"github.com/Carina-labs/HAL9000/client/base"
	"github.com/Carina-labs/HAL9000/client/base/query"
	novaTx "github.com/Carina-labs/HAL9000/client/nova/msgs"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
	"log"
	"os"
	"reflect"
	"time"
)

var (
	Host   = &config.HostChainInfo{}
	tmpseq = uint64(1)
)

func UpdateChainState(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File, botch chan<- time.Time) {

	Host.Set(host)

	conn, err := grpc.Dial(
		Host.GrpcAddr,
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
		botTickLog("Oracle", int(intv)*i)

		delegatedToken, height, apphash := OracleInfo(cq, Host.Validator)
		msg1 := novaTx.MakeMsgUpdateChainState(botInfo.GetAddress(), host, Host.Denom, delegatedToken, height, apphash)
		//msg2, _ := commonTx.MakeMsgSend(botInfo.GetAddress(), "nova1z36nmc2efth7wy3dcnjsw2tu83qn5mxyydu663", []string{"unova"}, []int64{1000})
		msgs := []sdktypes.Msg{msg1}
		for {
			ok := base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
			if ok {
				break
			}
		}
		log.Println("----> MsgUpdateChainState was sent")
		botch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaAutoStake(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File, botch chan<- time.Time) {

	Host.Set(host)

	conn, err := grpc.Dial(
		Host.GrpcAddr,
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
		botTickLog("Re-Staking", int(intv)*i)

		r := RewardsWithAddr(cq, Host.HostAccount, Host.Validator)
		if reflect.DeepEqual(r, sdktypes.DecCoin{}) {
			time.Sleep(intv * time.Second)
			i++
			continue
		}

		msg1 := novaTx.MakeMsgIcaAutoStaking(host, botInfo.GetAddress(), r)
		msgs := []sdktypes.Msg{msg1}
		for {
			ok := base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
			if ok {
				break
			}
		}
		log.Println("----> MsgIcaAutoStaking was sent")
		botch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaStake(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File, botch chan<- time.Time) {

	i := 0
	intv := time.Duration(interval)
	for {
		botTickLog("ICA-Staking", int(intv)*i)

		msg1 := novaTx.MakeMsgDelegate(host, botInfo.GetAddress(), tmpseq)
		msgs := []sdktypes.Msg{msg1}
		for {
			ok := base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
			if ok {
				break
			}
		}
		log.Println("----> MsgDelegate was sent")
		botch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func UndelegateAndWithdraw(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, chanID string, interval int, errLogger *os.File, botch chan<- time.Time) {

	Host.Set(host)

	conn, err := grpc.Dial(
		Host.GrpcAddr,
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
		botTickLog("Undelegate & Withdraw", int(intv)*i)

		blkTS := LatestBlockTS(cq)
		delegatedToken, height, apphash := OracleInfo(cq, Host.Validator)
		msg1 := novaTx.MakeMsgUpdateChainState(botInfo.GetAddress(), host, Host.Denom, delegatedToken, height, apphash)

		if isStart {
			msg2 := novaTx.MakeMsgUndelegate(host, botInfo.GetAddress(), tmpseq)
			msgs := []sdktypes.Msg{msg1, msg2}
			for {
				ok := base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
				if ok {
					break
				}
			}
			log.Println("----> MsgUndelegate was sent")
			botch <- time.Now().UTC()
			isStart = false
		} else {
			msg2 := novaTx.MakeMsgUndelegate(host, botInfo.GetAddress(), tmpseq)
			msgs := []sdktypes.Msg{msg1, msg2}
			for {
				ok := base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
				if ok {
					break
				}
			}
			log.Println("----> MsgUndelegate was sent")

			time.Sleep(60 * time.Second)

			msg3 := novaTx.MakeMsgIcaWithdraw(host, botInfo.GetAddress(), "transfer", chanID, blkTS, tmpseq)
			msgs = []sdktypes.Msg{msg3}
			for {
				ok := base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
				if ok {
					break
				}
			}
			log.Println("----> MsgPendingWithdraw was sent")
			botch <- time.Now().UTC()
		}

		time.Sleep(intv * time.Second)
		i++
	}

}
