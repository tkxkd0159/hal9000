package main

import (
	"fmt"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"os"
	"path"
	"sync"
)

var wg sync.WaitGroup
var wd string
var botInfo keyring.Info

func init() {
	cwd, err := os.Getwd()
	utils.HandleErr(err, "cannot get working directory", ut.EXIT)
	wd = path.Join(cwd, "bot")
	err = os.Mkdir(wd, 0740)
	if os.IsExist(err) {
		log.Println("bot directory already exist")
	} else if err != nil {
		log.Fatal(err)
	}
}

func main() {
	sViper := config.Sviper
	ctx, _ := common.MakeContext(novaapp.ModuleBasics, viper.GetString("nova.local_addr"),
		"tcp://localhost:26657", viper.GetString("nova.chain_id"), wd, keyring.BackendFile)
	brandnew := false
	if brandnew {
		botInfo = common.MakeClientWithNewAcc(ctx, "nova-bot", sViper.GetString("nova_mne"), sdktypes.FullFundraiserPath, hd.Secp256k1)
	} else {
		botInfo = common.LoadClient(ctx, "nova-bot")
	}
	fmt.Println(botInfo.GetPubKey(), botInfo.GetName(), botInfo.GetType(), botInfo.GetAddress(), botInfo.GetAlgo())
	fmt.Println(botInfo.GetPath())

	addr := viper.GetString("net.ip.nova") + ":" + viper.GetString("net.port.grpc")
	fmt.Println(addr)
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
