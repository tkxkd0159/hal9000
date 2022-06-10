package main

import (
	"context"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)
import (
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	stakingv1beta1 "github.com/Carina-labs/nova/api/cosmos/staking/v1beta1"
)

func main() {
	nv := config.SetNetInfo()
	addr := viper.GetString("NOVA_ADDR")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.HandleErr(err)
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	go func() {
		c := tendermintv1beta1.NewServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.GetNodeInfo(ctx, &tendermintv1beta1.GetNodeInfoRequest{})
		utils.HandleErr(err)
		//log.Printf("From gRPC srv : \n %+v", r.GetApplicationVersion())
		log.Printf("From gRPC srv : \n %+v", r.GetNodeInfo())
	}()

	go func() {
		c := stakingv1beta1.NewQueryClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.Validator(ctx, &stakingv1beta1.QueryValidatorRequest{ValidatorAddr: nv.GetString("nova.val_addr")})
		utils.HandleErr(err)
		log.Printf("From gRPC srv : \n %+v", r.GetValidator().Tokens)
	}()
	time.Sleep(20 * time.Second)

}
