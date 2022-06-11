package common

import (
	"context"
	"fmt"
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"
	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	stakingv1beta1 "github.com/Carina-labs/nova/api/cosmos/staking/v1beta1"
	"google.golang.org/grpc"
	"time"
)

func GetNodeInfo(conn *grpc.ClientConn) *tendermintv1beta1.GetNodeInfoResponse {
	c := tendermintv1beta1.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	r, err := c.GetNodeInfo(ctx, &tendermintv1beta1.GetNodeInfoRequest{})
	utils.HandleErr(err, "Can't get Node info", types.KEEP)
	// log.Printf("From gRPC srv : \n %+v", r.GetApplicationVersion())
	//log.Printf("From gRPC srv : \n %+v", r.GetNodeInfo())
	return r
}

func GetValInfo(conn *grpc.ClientConn, valAddr string) *stakingv1beta1.QueryValidatorResponse {
	c := stakingv1beta1.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Validator(ctx, &stakingv1beta1.QueryValidatorRequest{ValidatorAddr: valAddr})
	utils.HandleErr(err, fmt.Sprintf("Can't get %s info", valAddr), types.KEEP)
	// r.GetValidator().Tokens
	return r
}

// SendTx
// 1. Generate a TX with Msg (TxBuilder)
// 2. Sign the generated transaction with the keyring's account
// 3. Broadcast the tx using gPRC
func SendTx() {

}
