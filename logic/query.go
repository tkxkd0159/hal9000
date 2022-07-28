package logic

import (
	"github.com/Carina-labs/HAL9000/client/common/query"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func OracleInfo(cq *query.CosmosQueryClient, validatorAddr string) (string, int64, []byte) {
	h := cq.GetLatestBlock().GetBlock().GetHeader().GetHeight()
	hisInfo := cq.GetHistoricalInfo(h)
	apphash := hisInfo.GetHist().GetHeader().GetAppHash()

	var delegatedToken string
	for _, val := range hisInfo.GetHist().GetValset() {
		if val.GetOperatorAddress() != validatorAddr {
			continue
		} else {
			delegatedToken = val.GetTokens()
			break
		}
	}
	return delegatedToken, h, apphash
}

func LatestBlockTS(cq *query.CosmosQueryClient) time.Time {
	secs := cq.GetLatestBlock().GetBlock().GetHeader().GetTime().GetSeconds()
	nanos := cq.GetLatestBlock().GetBlock().GetHeader().GetTime().GetNanos()
	currentTs := time.Unix(secs, int64(nanos)).UTC()
	return currentTs
}

func RewardsWithAddr(cq *query.CosmosQueryClient, delegator string, validator string) sdktypes.DecCoin {
	return cq.GetRewards(delegator, validator).GetRewards()[0]
}