package main

import (
	"flag"
	"fmt"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	cq "github.com/Carina-labs/HAL9000/client/common/query"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
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

var (
	wg sync.WaitGroup
)

var (
	krDir, logDir string
	ctx           client.Context
	botInfo       keyring.Info
	sViper        *viper.Viper
)

func init() {
	sViper = config.Sviper
	common.SetBechPrefix()
	krDir, logDir = SetInitialDir("bot", "logs")
}

// FIXME: wasmvm doesn't support AArch64. Need to set GOARCH=amd64
func main() {
	keyname := flag.String("name", "nova-bot", "unique key name")
	newacc := flag.Bool("add", false, "Start client with making new account")
	disp := flag.Bool("display", false, "show context log through stdout")

	userOutput := os.Stdout
	flag.Parse()
	if !*disp {
		fpLog, err := os.OpenFile(path.Join(logDir, "ctxlog.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "", 0)
		userOutput = fpLog
	}
	defer userOutput.Close()

	novaGrpcAddr := viper.GetString("net.ip.nova") + ":" + viper.GetString("net.port.grpc")
	novaTmAddr := "tcp://" + viper.GetString("net.ip.nova") + ":" + viper.GetString("net.port.tmrpc")

	pp := GetPassphrase(sViper)

	// set pipe to ignore stdin tty
	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)

	if *newacc {

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			viper.GetString("nova.local_addr"),
			novaTmAddr,
			viper.GetString("nova.chain_id"),
			krDir,
			keyring.BackendFile,
			os.Stdin,
			userOutput,
		)

		botInfo = common.MakeClientWithNewAcc(
			ctx,
			*keyname,
			sViper.GetString("nova_mne"),
			sdktypes.FullFundraiserPath,
			hd.Secp256k1,
		)
		os.Exit(0)
	} else {
		_, err = wpipe.Write([]byte(pp))
		utils.CheckErr(err, "", 0)

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			viper.GetString("nova.local_addr"),
			novaTmAddr,
			viper.GetString("nova.chain_id"),
			krDir,
			keyring.BackendFile,
			rpipe,
			userOutput,
		)
		os.Stdin = rpipe

		botInfo = common.LoadClientPubInfo(ctx, *keyname)
	}
	ctx = common.AddMoreFromInfo(ctx)
	txf := common.MakeTxFactory(ctx, "auto", "0unova", "", 1.1)

	// ** Build TX
	msg1, _ := common.MakeSendMsg(botInfo.GetAddress(), viper.GetString("nova.target_addr"), []string{"unova"}, []int64{777})
	msgs := []sdktypes.Msg{msg1}
	common.GenTxWithFactory(ctx, txf, false, msgs...)

	wg.Add(2)
	go func() {
		defer wg.Done()
		api.Server{}.On("127.0.0.1:3334")
	}()

	conn, err := grpc.Dial(
		novaGrpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", 0)

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("unexpected gRPC disconnection: %v", err)
		}
	}(conn)

	fmt.Println("\n************ gRPC query checking ************")
	nf := cq.GetNodeRes(conn)
	vf := cq.GetValRes(conn, viper.GetString("nova.val_addr"))
	fmt.Println(nf.GetNodeInfo())
	fmt.Println(vf.GetValidator().Tokens)
	wg.Wait()
}
