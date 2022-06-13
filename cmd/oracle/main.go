package main

import (
	"flag"
	"fmt"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/client/common/types"
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
	wg        sync.WaitGroup
	origStdin *os.File
)

var (
	wd      string
	ctx     client.Context
	botInfo keyring.Info
	sViper  *viper.Viper
)

func init() {
	cwd, err := os.Getwd()
	utils.CheckErr(err, "cannot get working directory", 0)
	wd = path.Join(cwd, "bot")
	err = os.Mkdir(wd, 0740)
	if os.IsExist(err) {
		log.Println("** bot directory already exist **")
	} else if err != nil {
		log.Fatal(err)
	}

	sViper = config.Sviper

	config := sdktypes.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	config.Seal()
}

// FIXME: wasmvm doesn't support AArch64. Need to set GOARCH=amd64
// make run TARGET=oracle CUSTOM_ORGS="-add=true -name='gogo'"
func main() {
	keyname := flag.String("name", "nova-bot", "unique key name")
	newacc := flag.Bool("add", false, "Start client with making new account")
	flag.Parse()
	novaGrpcAddr := viper.GetString("net.ip.nova") + ":" + viper.GetString("net.port.grpc")
	novaTmAddr := "tcp://" + viper.GetString("net.ip.nova") + ":" + viper.GetString("net.port.tmrpc")

	mypw := sViper.GetString("pw")
	mypphrase := fmt.Sprintf("%s\n%s\n", mypw, mypw)
	// set pipe to ignore stdin tty
	rpipe, wpipe, err := utils.SetPipe(origStdin)
	utils.CheckErr(err, "", 0)

	if *newacc {
		ctx, _ = common.MakeContext(novaapp.ModuleBasics, viper.GetString("nova.local_addr"),
			novaTmAddr, viper.GetString("nova.chain_id"), wd, keyring.BackendFile, os.Stdin)
		botInfo = common.MakeClientWithNewAcc(ctx, *keyname, sViper.GetString("nova_mne"), sdktypes.FullFundraiserPath, hd.Secp256k1)
		os.Exit(0)
	} else {
		_, err = wpipe.Write([]byte(mypphrase))
		utils.CheckErr(err, "", 0)

		ctx, _ = common.MakeContext(novaapp.ModuleBasics, viper.GetString("nova.local_addr"),
			novaTmAddr, viper.GetString("nova.chain_id"), wd, keyring.BackendFile, rpipe)
		os.Stdin = rpipe

		botInfo = common.LoadClient(ctx, *keyname)
	}

	ctx = common.AddMoreFromInfo(ctx)

	// ** Build TX
	txf := common.MakeTxFactory(ctx, "auto", "0unova", "0unova", "")
	coin := sdktypes.Coin{Denom: "unova", Amount: sdktypes.NewInt(333)}
	to := sdktypes.AccAddress(viper.GetString("nova.target_addr"))
	msg1 := common.MakeSendMsg(botInfo.GetAddress(), to, coin)
	msgs := []sdktypes.Msg{msg1}

	common.GenTxWithFactory(ctx, txf, false, msgs...)

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

	wg.Add(2)

	go func() {
		defer wg.Done()
		api.Server{}.On()
	}()

	fmt.Println("\n************ gRPC query checking ************\n")
	nf := common.GetNodeInfo(conn)
	vf := common.GetValInfo(conn, viper.GetString("nova.val_addr"))

	fmt.Println(nf.GetNodeInfo())
	fmt.Println(vf.GetValidator().Tokens)
	wg.Wait()
}
