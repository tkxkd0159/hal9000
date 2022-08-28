package query

import (
	"context"
	"fmt"
	"github.com/Carina-labs/HAL9000/client/base/types"
	"github.com/Carina-labs/HAL9000/utils"
	utiltypes "github.com/Carina-labs/HAL9000/utils/types"
	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	bankv1beta1 "github.com/cosmos/cosmos-sdk/x/bank/types"
	distv1beta1 "github.com/cosmos/cosmos-sdk/x/distribution/types"

	stakingv1beta1 "github.com/Carina-labs/nova/api/cosmos/staking/v1beta1"
	txv1beta1 "github.com/cosmos/cosmos-sdk/types/tx"
	"google.golang.org/grpc"
	"time"
)

type CosmosQueryClient struct {
	*grpc.ClientConn
}

var (
	_ types.BaseQuerier = &CosmosQueryClient{}
)

const (
	ctxTimeout = time.Second * 5
)

// ######################### Tendermint #########################

func (cqc *CosmosQueryClient) GetNodeRes() *tendermintv1beta1.GetNodeInfoResponse {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.GetNodeInfo(ctx, &tendermintv1beta1.GetNodeInfoRequest{})
	utils.CheckErr(err, "Can't get Node info", utiltypes.KEEP)

	return r
}

func (cqc *CosmosQueryClient) GetLatestBlock() *tendermintv1beta1.GetLatestBlockResponse {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.GetLatestBlock(ctx, &tendermintv1beta1.GetLatestBlockRequest{})
	utils.CheckErr(err, "", 1)
	return r
}

func (cqc *CosmosQueryClient) GetBlockByHeight(height int64) *tendermintv1beta1.GetBlockByHeightResponse {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.GetBlockByHeight(ctx, &tendermintv1beta1.GetBlockByHeightRequest{Height: height})
	utils.CheckErr(err, "", 1)
	return r
}

// ######################### Bank #########################

func (cqc *CosmosQueryClient) GetBalance(addr string, denom string) *bankv1beta1.QueryBalanceResponse {
	c := bankv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.Balance(ctx, &bankv1beta1.QueryBalanceRequest{Address: addr, Denom: denom})
	utils.CheckErr(err, "", 1)
	return r
}

// ######################### Staking #########################

func (cqc *CosmosQueryClient) GetValInfo(valAddr string) *stakingv1beta1.QueryValidatorResponse {
	c := stakingv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.Validator(ctx, &stakingv1beta1.QueryValidatorRequest{ValidatorAddr: valAddr})
	utils.CheckErr(err, fmt.Sprintf("Can't get %s info", valAddr), utiltypes.KEEP)
	return r
}

func (cqc *CosmosQueryClient) GetHistoricalInfo(height int64) *stakingv1beta1.QueryHistoricalInfoResponse {
	c := stakingv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.HistoricalInfo(ctx, &stakingv1beta1.QueryHistoricalInfoRequest{Height: height})
	utils.CheckErr(err, "", 1)
	return r
}

// ######################### Tx #########################

func (cqc *CosmosQueryClient) GetTx(hash string) *txv1beta1.GetTxResponse {
	c := txv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.GetTx(ctx, &txv1beta1.GetTxRequest{Hash: hash})
	utils.CheckErr(err, "", 1)
	return r
}

func (cqc *CosmosQueryClient) GetRewards(delegator string, validator string) *distv1beta1.QueryDelegationRewardsResponse {
	c := distv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.DelegationRewards(ctx, &distv1beta1.QueryDelegationRewardsRequest{DelegatorAddress: delegator, ValidatorAddress: validator})
	utils.CheckErr(err, "", 1)
	return r
}
