package common

import (
	"context"
	"github.com/Carina-labs/HAL9000/client"
	"github.com/Carina-labs/HAL9000/utils"
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
	utils.HandleErr(err)
	// log.Printf("From gRPC srv : \n %+v", r.GetApplicationVersion())
	//log.Printf("From gRPC srv : \n %+v", r.GetNodeInfo())
	return r
}

func GetValInfo(conn *grpc.ClientConn) *stakingv1beta1.QueryValidatorResponse {
	c := stakingv1beta1.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Validator(ctx, &stakingv1beta1.QueryValidatorRequest{ValidatorAddr: client.NV.GetString("nova.val_addr")})
	utils.HandleErr(err)
	// r.GetValidator().Tokens
	return r
}
