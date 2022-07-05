package query

import (
	"context"
	"fmt"
	"github.com/Carina-labs/HAL9000/client/common/types"
	"github.com/Carina-labs/HAL9000/utils"
	utiltypes "github.com/Carina-labs/HAL9000/utils/types"
	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	stakingv1beta1 "github.com/Carina-labs/nova/api/cosmos/staking/v1beta1"
	txv1beta1 "github.com/cosmos/cosmos-sdk/types/tx"
	"google.golang.org/grpc"
	"time"
)

type CosmosQueryClient struct {
	*grpc.ClientConn
	A string
}

var (
	_ types.CommonQuery = &CosmosQueryClient{}
)

// ######################### Tendermint #########################

func (cqc *CosmosQueryClient) GetNodeRes() *types.NodeInfoRes {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetNodeInfo(ctx, &types.NodeInfoReq{})
	utils.CheckErr(err, "Can't get Node info", utiltypes.KEEP)

	return r
}

func (cqc *CosmosQueryClient) GetLatestBlock() *tendermintv1beta1.GetLatestBlockResponse {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetLatestBlock(ctx, &tendermintv1beta1.GetLatestBlockRequest{})
	utils.CheckErr(err, "", 1)
	return r
}

func (cqc *CosmosQueryClient) GetBlockByHeight(height int64) *tendermintv1beta1.GetBlockByHeightResponse {
	c := tendermintv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetBlockByHeight(ctx, &tendermintv1beta1.GetBlockByHeightRequest{Height: height})
	utils.CheckErr(err, "", 1)
	return r
}

// ######################### Staking #########################

func (cqc *CosmosQueryClient) GetValInfo(valAddr string) *types.ValInfoRes {
	c := stakingv1beta1.NewQueryClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Validator(ctx, &types.ValInfoReq{ValidatorAddr: valAddr})
	utils.CheckErr(err, fmt.Sprintf("Can't get %s info", valAddr), utiltypes.KEEP)
	return r
}

// ######################### Tx #########################

func (cqc *CosmosQueryClient) GetTx(hash string) *txv1beta1.GetTxResponse {
	c := txv1beta1.NewServiceClient(cqc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetTx(ctx, &txv1beta1.GetTxRequest{Hash: hash})
	utils.CheckErr(err, "", 1)
	return r
}
