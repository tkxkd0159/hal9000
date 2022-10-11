package logic

import (
	"reflect"
	"sync"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/Carina-labs/HAL9000/client/base"
	"github.com/Carina-labs/HAL9000/client/base/query"
	novatypes "github.com/Carina-labs/HAL9000/client/base/types"
	novaTx "github.com/Carina-labs/HAL9000/client/nova/msgs"
	novaq "github.com/Carina-labs/HAL9000/client/nova/query"
	"github.com/Carina-labs/HAL9000/config"
)

func mustExecTx(b *novatypes.Bot, host *config.HostChainInfo, msgs []sdktypes.Msg, opts ...IBCConfirm) {
	if opts != nil {
		ibc := opts[0]
	LOOP1:
		for {
			var done bool
			switch base.GenTxByBot(b, msgs...) {
			case base.NEXT:
				break LOOP1
			case base.NONE:
				time.Sleep(IBCDelay)
				done = isIBCDone(ibc.seq, FetchBotSeq(ibc.nq, ibc.action, host.Name))
			case base.CRITICAL:
				done = false
			}
			if done {
				botMsgLog(msgs)
				break
			}
		}
	} else {
	LOOP2:
		for {
			var done bool
			switch base.GenTxByBot(b, msgs...) {
			case base.NEXT:
				break LOOP2
			case base.NONE:
				done = true
			case base.CRITICAL:
				done = false
			}
			if done {
				botMsgLog(msgs)
				break
			}
		}
	}
}

func UpdateChainState(cq *query.CosmosQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Oracle", int(intv)*i, b.Interval)

		delegatedToken, height, apphash := OracleInfo(cq, host.Validator, host.HostAccount)
		msg1 := novaTx.MakeMsgUpdateChainState(b.KrInfo.GetAddress(), host.Name, host.Denom, delegatedToken, height, apphash)
		//msg2, _ := commonTx.MakeMsgSend(botInfo.GetAddress(), "nova1z36nmc2efth7wy3dcnjsw2tu83qn5mxyydu663", []string{"unova"}, []int64{1000})
		msgs := []sdktypes.Msg{msg1}
		mustExecTx(b, host, msgs)
		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaAutoStake(cq *query.CosmosQueryClient, nq *novaq.NovaQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Re-Staking", int(intv)*i, b.Interval)

		hostReward := RewardsWithAddr(cq, host.HostAccount, host.Validator)
		if reflect.DeepEqual(hostReward, sdktypes.DecCoin{}) {
			time.Sleep(intv * time.Second)
			i++
			continue
		}

		targetSeq := FetchBotSeq(nq, config.ActAutoStake, host.Name)
		msg1 := novaTx.MakeMsgIcaAutoStaking(host.Name, b.KrInfo.GetAddress(), hostReward, targetSeq)
		msgs := []sdktypes.Msg{msg1}
		mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActAutoStake, targetSeq})
		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaStake(nq *novaq.NovaQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("ICA-Staking", int(intv)*i, b.Interval)

		targetSeq := FetchBotSeq(nq, config.ActStake, host.Name)
		msg1 := novaTx.MakeMsgDelegate(host.Name, b.KrInfo.GetAddress(), targetSeq)
		msgs := []sdktypes.Msg{msg1}
		mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActStake, targetSeq})

		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func UndelegateAndWithdraw(cq *query.CosmosQueryClient, nq *novaq.NovaQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	isStart := true
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Undelegate & Withdraw", int(intv)*i, b.Interval)

		delegatedToken, height, apphash := OracleInfo(cq, host.Validator, host.HostAccount)
		msg1 := novaTx.MakeMsgUpdateChainState(b.KrInfo.GetAddress(), host.Name, host.Denom, delegatedToken, height, apphash)

		if isStart {
			undelSeq := FetchBotSeq(nq, config.ActUndelegate, host.Name)
			msg2 := novaTx.MakeMsgUndelegate(host.Name, b.KrInfo.GetAddress(), undelSeq)
			msgs := []sdktypes.Msg{msg1, msg2}
			mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActUndelegate, undelSeq})
			b.APIch <- time.Now().UTC()
			isStart = false
		} else {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				wg.Done()

				undelSeq := FetchBotSeq(nq, config.ActUndelegate, host.Name)
				msg2 := novaTx.MakeMsgUndelegate(host.Name, b.KrInfo.GetAddress(), undelSeq)
				msgs := []sdktypes.Msg{msg1, msg2}
				mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActUndelegate, undelSeq})
			}()

			go func() {
				wg.Done()

				wdSeq := FetchBotSeq(nq, config.ActWithdraw, host.Name)
				blkTS := LatestBlockTS(cq)
				msg3 := novaTx.MakeMsgIcaWithdraw(host.Name, b.KrInfo.GetAddress(), "transfer", host.IBCInfo.Transfer, blkTS, wdSeq)
				msgs := []sdktypes.Msg{msg3}
				mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActWithdraw, wdSeq})
			}()

			wg.Wait()
		}
		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}

}
