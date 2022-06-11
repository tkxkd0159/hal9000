package main

import (
	"fmt"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client"
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
)

func init() {
	config.SetEnv()
}

var wg sync.WaitGroup

func main() {
	addr := viper.GetString("NOVA_ADDR")
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	utils.HandleErr(err, "cannot create gRPC connection", ut.EXIT)

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("unexpected gRPC disconnection: %v", err)
		}
	}(conn)

	wg.Add(2)

	go func() {
		defer wg.Done()
		api.Server{}.On()
	}()

	ch1 := make(chan string)
	go func() {
		defer wg.Done()
		ch1 <- "get ch1"
	}()

	nf := common.GetNodeInfo(conn)
	vf := common.GetValInfo(conn, client.NV.GetString("nova.val_addr"))

	fmt.Println(nf.GetNodeInfo())
	fmt.Println(vf.GetValidator().Tokens)

	fmt.Println(<-ch1)
	wg.Wait()

}
