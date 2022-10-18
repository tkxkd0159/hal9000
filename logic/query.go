package logic

import (
	"fmt"

	"log"
	"time"

	tmtypes "github.com/Carina-labs/nova/api/tendermint/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Carina-labs/HAL9000/client/base/query"
	nquery "github.com/Carina-labs/HAL9000/client/nova/query"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
)

func OracleInfo(cq *query.CosmosQueryClient, validatorAddr, delegatorAddr string) (string, int64, []byte) {
	var bh int64
	var apphash []byte
	var delegatedToken string

	for {
		res, err := cq.GetLatestBlock()
		if err == nil {
			bh = res.GetBlock().GetHeader().GetHeight()
			break
		} else {
			utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
		}
		time.Sleep(ReQueryDelay)
	}

	for {
		res, err := cq.GetHistoricalInfo(bh)
		if err == nil {
			apphash = res.GetHist().GetHeader().GetAppHash()
			break
		} else {
			utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
		}
		time.Sleep(ReQueryDelay)
	}

	for {
		res, err := cq.GetDelegation(validatorAddr, delegatorAddr)
		if err == nil {
			delegatedToken = res.GetDelegationResponse().GetBalance().GetAmount()
			break
		} else {
			if err.Error() == status.Errorf(
				codes.NotFound,
				"delegation with delegator %s not found for validator %s",
				delegatorAddr, validatorAddr).Error() {
				delegatedToken = "0"
				break
			} else {
				utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
			}
		}
		time.Sleep(ReQueryDelay)
	}

	fmt.Println(" 🤪 Current delegated token on this target : ", delegatedToken)

	return delegatedToken, bh, apphash
}

func LatestBlockTS(cq *query.CosmosQueryClient) (ts time.Time) {
	var blk *tmtypes.Block
	for {
		res, err := cq.GetLatestBlock()
		if err == nil {
			blk = res.GetBlock()
			break
		} else {
			utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
		}
		time.Sleep(ReQueryDelay)
	}
	secs := blk.GetHeader().GetTime().GetSeconds()
	nanos := blk.GetHeader().GetTime().GetNanos()
	currentTs := time.Unix(secs, int64(nanos)).UTC()
	return currentTs
}

func RewardsWithAddr(cq *query.CosmosQueryClient, delegator string, validator string) (reward sdktypes.DecCoin) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("There is no reward to handle")
			reward = sdktypes.DecCoin{}
		}
	}()
	for {
		res, err := cq.GetRewards(delegator, validator)
		if err == nil {
			reward = res.GetRewards()[0]
			break
		} else {
			utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
		}
		time.Sleep(ReQueryDelay)
	}
	return
}

func FetchBotSeq(nq *nquery.NovaQueryClient, action string, zoneid string) (seq uint64) {
LOOP:
	for {
		switch action {
		case config.ActStake:
			res, err := nq.CurrentDelegateVersion(zoneid)
			if err == nil {
				seq = res.GetVersion()
				break LOOP
			} else {
				utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
			}
		case config.ActAutoStake:
			res, err := nq.CurrentAutoStakingVersion(zoneid)
			if err == nil {
				seq = res.GetVersion()
				break LOOP
			} else {
				utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
			}
		case config.ActUndelegate:
			res, err := nq.CurrentUndelegateVersion(zoneid)
			if err == nil {
				seq = res.GetVersion()
				break LOOP
			} else {
				utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
			}
		case config.ActWithdraw:
			res, err := nq.CurrentWithdrawVersion(zoneid)
			if err == nil {
				seq = res.GetVersion()
				break LOOP
			} else {
				utils.CheckErr(err, QueryErrPrefix, ut.KEEP)
			}
		default:
			panic("there is no sequence on this action")
		}

		time.Sleep(ReQueryDelay)
	}
	return
}

type IBCConfirm struct {
	nq     *nquery.NovaQueryClient
	action string
	seq    uint64
}

func isIBCDone(before, after uint64) bool {
	return (before + 1) == after
}
