package logic

import (
	"reflect"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/Carina-labs/HAL9000/client/base"
	"github.com/Carina-labs/HAL9000/client/base/query"
	novatypes "github.com/Carina-labs/HAL9000/client/base/types"
	novaTx "github.com/Carina-labs/HAL9000/client/nova/msgs"
	"github.com/Carina-labs/HAL9000/config"
)

var (
	tmpseq = uint64(1)
)

func UpdateChainState(cq *query.CosmosQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Oracle", int(intv)*i)

		delegatedToken, height, apphash := OracleInfo(cq, host.Validator)
		msg1 := novaTx.MakeMsgUpdateChainState(b.KrInfo.GetAddress(), host.Name, host.Denom, delegatedToken, height, apphash)
		//msg2, _ := commonTx.MakeMsgSend(botInfo.GetAddress(), "nova1z36nmc2efth7wy3dcnjsw2tu83qn5mxyydu663", []string{"unova"}, []int64{1000})
		msgs := []sdktypes.Msg{msg1}
		for {
			ok := base.GenTxByBot(b, false, msgs...)
			if ok {
				break
			}
		}
		botMsgLog(msgs)
		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaAutoStake(cq *query.CosmosQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Re-Staking", int(intv)*i)

		r := RewardsWithAddr(cq, host.HostAccount, host.Validator)
		if reflect.DeepEqual(r, sdktypes.DecCoin{}) {
			time.Sleep(intv * time.Second)
			i++
			continue
		}

		msg1 := novaTx.MakeMsgIcaAutoStaking(host.Name, b.KrInfo.GetAddress(), r)
		msgs := []sdktypes.Msg{msg1}
		for {
			ok := base.GenTxByBot(b, false, msgs...)
			if ok {
				break
			}
		}
		botMsgLog(msgs)
		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaStake(b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("ICA-Staking", int(intv)*i)

		msg1 := novaTx.MakeMsgDelegate(host.Name, b.KrInfo.GetAddress(), tmpseq)
		msgs := []sdktypes.Msg{msg1}
		for {
			ok := base.GenTxByBot(b, false, msgs...)
			if ok {
				break
			}
		}
		botMsgLog(msgs)
		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func UndelegateAndWithdraw(cq *query.CosmosQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	isStart := true
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Undelegate & Withdraw", int(intv)*i)

		blkTS := LatestBlockTS(cq)
		delegatedToken, height, apphash := OracleInfo(cq, host.Validator)
		msg1 := novaTx.MakeMsgUpdateChainState(b.KrInfo.GetAddress(), host.Name, host.Denom, delegatedToken, height, apphash)

		if isStart {
			msg2 := novaTx.MakeMsgUndelegate(host.Name, b.KrInfo.GetAddress(), tmpseq)
			msgs := []sdktypes.Msg{msg1, msg2}
			for {
				ok := base.GenTxByBot(b, false, msgs...)
				if ok {
					break
				}
			}
			botMsgLog(msgs)
			b.APIch <- time.Now().UTC()
			isStart = false
		} else {
			msg2 := novaTx.MakeMsgUndelegate(host.Name, b.KrInfo.GetAddress(), tmpseq)
			msgs := []sdktypes.Msg{msg1, msg2}
			for {
				ok := base.GenTxByBot(b, false, msgs...)
				if ok {
					break
				}
			}
			botMsgLog(msgs)

			time.Sleep(60 * time.Second)

			msg3 := novaTx.MakeMsgIcaWithdraw(host.Name, b.KrInfo.GetAddress(), "transfer", host.IBCInfo.Transfer, blkTS, tmpseq)
			msgs = []sdktypes.Msg{msg3}
			for {
				ok := base.GenTxByBot(b, false, msgs...)
				if ok {
					break
				}
			}
			botMsgLog(msgs)
			b.APIch <- time.Now().UTC()
		}

		time.Sleep(intv * time.Second)
		i++
	}

}
