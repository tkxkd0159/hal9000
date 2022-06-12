package main

import (
	"fmt"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {
	sViper := config.Sviper
	fmt.Println(sViper.Get("atom_mne"))
	addr := viper.GetString("net.ip.nova") + ":" + viper.GetString("net.port.grpc")
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
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

	nf := common.GetNodeInfo(conn)
	vf := common.GetValInfo(conn, viper.GetString("nova.val_addr"))

	fmt.Println(nf.GetNodeInfo())
	fmt.Println(vf.GetValidator().Tokens)
	wg.Wait()
}
