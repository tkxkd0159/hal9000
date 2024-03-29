package logic

import (
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/tkxkd0159/HAL9000/client/base"
	basem "github.com/tkxkd0159/HAL9000/client/base/msgs"
	"github.com/tkxkd0159/HAL9000/client/base/query"
	novatypes "github.com/tkxkd0159/HAL9000/client/base/types"
	novaTx "github.com/tkxkd0159/HAL9000/client/nova/msgs"
	novaq "github.com/tkxkd0159/HAL9000/client/nova/query"
	"github.com/tkxkd0159/HAL9000/config"
	"github.com/tkxkd0159/HAL9000/utils"
	ut "github.com/tkxkd0159/HAL9000/utils/types"
)

func mustExecTx(b *novatypes.Bot, host *config.HostChainInfo, msgs []sdktypes.Msg, opts ...IBCConfirm) bool {
	if err := basem.Validate(msgs...); err != nil {
		utils.LogErrWithFd(b.ErrLogger, err, "[msg]", ut.KEEP)
		return false
	}

	if opts != nil {
		ibc := opts[0]
	IBCTX:
		for {
			var done bool
			txErr := base.GenTxByBot(b, msgs...)
			switch txErr {
			case base.NEXT:
				break IBCTX
			case base.NONE:
				time.Sleep(time.Duration(host.IBCTimeout) * time.Minute)
				done = isIBCDone(ibc.seq, FetchBotSeq(ibc.nq, ibc.action, host.Name))
			case base.CRITICAL:
				done = false
			case base.REPEAT:
				return false
			}

			if done {
				botMsgLog(msgs)
				break
			} else if txErr == base.NONE {
				// TODO: Need to implement a system that can resolve this issue more effectively than simply showing the log
				log.Printf(" 🚫 Caution : IBC Ack did not arrive normally during the set ibc timeout period ❗\n 🚫 Check the relayer status manually right now ❗\n️")
				os.Exit(1)
			}
		}
	} else {
	NORMTX:
		for {
			var done bool
			switch base.GenTxByBot(b, msgs...) {
			case base.NEXT:
				break NORMTX
			case base.NONE:
				done = true
			case base.CRITICAL:
				done = false
			case base.REPEAT:
				return false
			}
			if done {
				botMsgLog(msgs)
				break
			}
		}
	}

	return true
}

func UpdateChainState(cq *query.CosmosQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Oracle", int(intv)*i, b.Interval)
	ORACLE:
		delegatedToken, height, apphash := OracleInfo(cq, host.Validator, host.HostAccount)
		msg1 := novaTx.MakeMsgUpdateChainState(b.KrInfo.GetAddress(), host.Name, host.Denom, delegatedToken, height, apphash)
		msgs := []sdktypes.Msg{msg1}
		if ok := mustExecTx(b, host, msgs); !ok {
			goto ORACLE
		}
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
	RESTAKE:
		hostReward := RewardsWithAddr(cq, host.HostAccount, host.Validator)
		if reflect.DeepEqual(hostReward, sdktypes.DecCoin{}) {
			time.Sleep(intv * time.Second)
			i++
			continue
		}

		targetSeq := FetchBotSeq(nq, config.ActAutoStake, host.Name)
		msg1 := novaTx.MakeMsgIcaAutoStaking(host.Name, b.KrInfo.GetAddress(), hostReward, targetSeq, host.IBCTimeout)
		msgs := []sdktypes.Msg{msg1}
		if ok := mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActAutoStake, targetSeq}); !ok {
			goto RESTAKE
		}
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
	STAKE:
		targetSeq := FetchBotSeq(nq, config.ActStake, host.Name)
		msg1 := novaTx.MakeMsgDelegate(host.Name, b.KrInfo.GetAddress(), targetSeq, host.IBCTimeout)
		msgs := []sdktypes.Msg{msg1}
		if ok := mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActStake, targetSeq}); !ok {
			goto STAKE
		}

		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func UndelegateAndWithdraw(cq *query.CosmosQueryClient, nq *novaq.NovaQueryClient, b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Undelegate & Withdraw", int(intv)*i, b.Interval)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
		UNDEL:
			delegatedToken, height, apphash := OracleInfo(cq, host.Validator, host.HostAccount)
			undelSeq := FetchBotSeq(nq, config.ActUndelegate, host.Name)
			msg1 := novaTx.MakeMsgUpdateChainState(b.KrInfo.GetAddress(), host.Name, host.Denom, delegatedToken, height, apphash)
			msg2 := novaTx.MakeMsgUndelegate(host.Name, b.KrInfo.GetAddress(), undelSeq, host.IBCTimeout)
			msgs := []sdktypes.Msg{msg1, msg2}
			if ok := mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActUndelegate, undelSeq}); !ok {
				goto UNDEL
			}

			wg.Done()
		}()

		go func() {
		WD:
			wdSeq := FetchBotSeq(nq, config.ActWithdraw, host.Name)
			blkTS := LatestBlockTS(cq)
			msg3 := novaTx.MakeMsgIcaWithdraw(host.Name, b.KrInfo.GetAddress(), "transfer", host.IBCInfo.Transfer, blkTS, wdSeq, host.IBCTimeout)
			msgs := []sdktypes.Msg{msg3}
			if ok := mustExecTx(b, host, msgs, IBCConfirm{nq, config.ActWithdraw, wdSeq}); !ok {
				goto WD
			}

			wg.Done()
		}()

		wg.Wait()

		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}

func ClaimAllSnAsset(b *novatypes.Bot, host *config.HostChainInfo) {
	i := 0
	intv := time.Duration(b.Interval)
	for {
		botTickLog("Auto-Claim", int(intv)*i, b.Interval)
	CLAIM:
		msg1 := novaTx.MakeMsgClaimAllSnAsset(host.Name, b.KrInfo.GetAddress())
		msgs := []sdktypes.Msg{msg1}
		if ok := mustExecTx(b, host, msgs); !ok {
			goto CLAIM
		}

		b.APIch <- time.Now().UTC()
		time.Sleep(intv * time.Second)
		i++
	}
}
