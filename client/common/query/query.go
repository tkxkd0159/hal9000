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

func GetNodeRes(conn *grpc.ClientConn) *types.NodeInfoRes {
	c := tendermintv1beta1.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	r, err := c.GetNodeInfo(ctx, &types.NodeInfoReq{})
	utils.CheckErr(err, "Can't get Node info", utiltypes.KEEP)

	return r
}

func GetValRes(conn *grpc.ClientConn, valAddr string) *types.ValInfoRes {
	c := stakingv1beta1.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Validator(ctx, &types.ValInfoReq{ValidatorAddr: valAddr})
	utils.CheckErr(err, fmt.Sprintf("Can't get %s info", valAddr), utiltypes.KEEP)
	return r
}

func GetTx(conn *grpc.ClientConn, hash string) *txv1beta1.GetTxResponse {
	c := txv1beta1.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetTx(ctx, &txv1beta1.GetTxRequest{Hash: hash})
	utils.CheckErr(err, "", 1)
	return r
}
