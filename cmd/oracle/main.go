package main

import (
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)
import (
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
)

func init() {
	config.SetEnv()
}

func main() {

	addr := viper.GetString("NOVA_ADDR")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.HandleErr(err)

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("unexpected gRPC disconnection: %v", err)
		}
	}(conn)

	nf := common.GetNodeInfo(conn)
	vf := common.GetValInfo(conn)

	fmt.Println(nf.GetNodeInfo())
	fmt.Println(vf.GetValidator().Tokens)
}
