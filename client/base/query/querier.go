package query

import (
	"context"
	"time"

	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	stakingv1beta1 "github.com/Carina-labs/nova/api/cosmos/staking/v1beta1"
	txv1beta1 "github.com/cosmos/cosmos-sdk/types/tx"
	bankv1beta1 "github.com/cosmos/cosmos-sdk/x/bank/types"
	distv1beta1 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"google.golang.org/grpc"

	"github.com/Carina-labs/HAL9000/client/base/types"
	"github.com/Carina-labs/HAL9000/utils"
	utiltypes "github.com/Carina-labs/HAL9000/utils/types"
)

const (
	ctxTimeout = time.Second * 10
)

var (
	_ types.BaseQuerier = &CosmosQueryClient{}
)

type CosmosQueryClient struct {
	*grpc.ClientConn
}

func NewCosmosQueryClient(grpcAddr string) *CosmosQueryClient {
	conn, err := grpc.Dial(
		grpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", utiltypes.EXIT)
	return &CosmosQueryClient{conn}
}

// ######################### Tendermint #########################

func (cqc *CosmosQueryClient) GetNodeRes() (*tendermintv1beta1.GetNodeInfoResponse, error) {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.GetNodeInfo(ctx, &tendermintv1beta1.GetNodeInfoRequest{})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (cqc *CosmosQueryClient) GetLatestBlock() (*tendermintv1beta1.GetLatestBlockResponse, error) {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.GetLatestBlock(ctx, &tendermintv1beta1.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (cqc *CosmosQueryClient) GetBlockByHeight(height int64) (*tendermintv1beta1.GetBlockByHeightResponse, error) {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.GetBlockByHeight(ctx, &tendermintv1beta1.GetBlockByHeightRequest{Height: height})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ######################### Bank #########################

func (cqc *CosmosQueryClient) GetBalance(addr string, denom string) (*bankv1beta1.QueryBalanceResponse, error) {
	c := bankv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.Balance(ctx, &bankv1beta1.QueryBalanceRequest{Address: addr, Denom: denom})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ######################### Staking #########################

func (cqc *CosmosQueryClient) GetValInfo(valAddr string) (*stakingv1beta1.QueryValidatorResponse, error) {
	c := stakingv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.Validator(ctx, &stakingv1beta1.QueryValidatorRequest{ValidatorAddr: valAddr})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (cqc *CosmosQueryClient) GetHistoricalInfo(height int64) (*stakingv1beta1.QueryHistoricalInfoResponse, error) {
	c := stakingv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.HistoricalInfo(ctx, &stakingv1beta1.QueryHistoricalInfoRequest{Height: height})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (cqc *CosmosQueryClient) GetDelegation(valAddr, delAddr string) (*stakingv1beta1.QueryDelegationResponse, error) {
	c := stakingv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.Delegation(ctx, &stakingv1beta1.QueryDelegationRequest{DelegatorAddr: delAddr, ValidatorAddr: valAddr})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ######################### Tx #########################

func (cqc *CosmosQueryClient) GetTx(hash string) (*txv1beta1.GetTxResponse, error) {
	c := txv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.GetTx(ctx, &txv1beta1.GetTxRequest{Hash: hash})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (cqc *CosmosQueryClient) GetRewards(delegator string, validator string) (*distv1beta1.QueryDelegationRewardsResponse, error) {
	c := distv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.DelegationRewards(ctx, &distv1beta1.QueryDelegationRewardsRequest{DelegatorAddress: delegator, ValidatorAddress: validator})
	if err != nil {
		return nil, err
	}
	return r, nil
}
