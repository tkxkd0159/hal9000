package main

import (
	"flag"
	"fmt"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	cq "github.com/Carina-labs/HAL9000/client/common/query"
	novac "github.com/Carina-labs/HAL9000/client/nova"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
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
	"time"
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

	apiAddr := flag.String("api", "127.0.0.1:3334", "Set bot api address")
	keyname := flag.String("name", "nova-bot", "Set unique key name (uid)")
	newacc := flag.Bool("add", false, "Start client with making new account")
	intv := flag.Int("interval", 0, "Oracle update interval")
	disp := flag.Bool("display", false, "Show context log through stdout")
	flag.Parse()

	// Open api endpoint to check bot
	wg.Add(2)
	go func() {
		defer wg.Done()
		api.Server{}.On(*apiAddr)
	}()

	// Start bot logic
	userOutput := os.Stdout
	var fpErr *os.File
	var fpErrNova *os.File
	if !*disp {
		fpLog, err := os.OpenFile(path.Join(logDir, "ctxlog.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open logfp", 0)

		// 외부 라이브러리에서 fmt.Fprintf(os.stderr)로 처리하는 애들 핸들링
		fpErr, err = os.OpenFile(path.Join(logDir, "other_err.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open otherErr", 0)

		// 반환되서 처리할 수 있는 에러 핸들링
		fpErrNova, err = os.OpenFile(path.Join(logDir, "nova_err.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "cannot open novaerr", 0)

		userOutput = fpLog
		os.Stderr = fpErr
	}

	projFps := []*os.File{userOutput, fpErr, fpErrNova}
	defer func(fps ...*os.File) {
		for _, fp := range fps {
			err := fp.Close()
			if err != nil {
				utils.CheckErr(err, "", 1)
			}
		}
	}(projFps...)

	novaGrpcAddr := viper.GetString("net.ip.nova") + ":" + viper.GetString("net.port.grpc")
	novaTmAddr := "tcp://" + viper.GetString("net.ip.nova") + ":" + viper.GetString("net.port.tmrpc")

	pp := GetPassphrase(sViper)

	// set pipe to ignore stdin tty
	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)

	if *newacc {

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			viper.GetString("nova.bot_addr"),
			novaTmAddr,
			viper.GetString("nova.chain_id"),
			krDir,
			keyring.BackendFile,
			os.Stdin,
			userOutput,
			false,
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
			viper.GetString("nova.bot_addr"),
			novaTmAddr,
			viper.GetString("nova.chain_id"),
			krDir,
			keyring.BackendFile,
			rpipe,
			userOutput,
			false,
		)
		os.Stdin = rpipe

		botInfo = common.LoadClientPubInfo(ctx, *keyname)
	}
	ctx = common.AddMoreFromInfo(ctx)
	txf := common.MakeTxFactory(ctx, "auto", "0unova", "", 1.1)
	time.Now()
	// ** Build TX
	go func(interval int) {
		stream := ut.Fstream{Err: fpErrNova}
		for {
			msg1, _ := common.MakeMsgSend(botInfo.GetAddress(), viper.GetString("nova.target_addr"), []string{"unova"}, []int64{777})

			msg2, err := novac.MakeMsgUpdateChainState(botInfo.GetAddress(), "uatom", 7654321, 6, 777)
			if err != nil {
				continue
			}
			msgs := []sdktypes.Msg{msg1, msg2}

			common.GenTxWithFactory(stream, ctx, txf, false, msgs...)
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}(*intv)

	conn, err := grpc.Dial(
		novaGrpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", 0)

	defer func(c *grpc.ClientConn) {
		err := c.Close()
		if err != nil {
			log.Printf("unexpected gRPC disconnection: %v", err)
		}
	}(conn)

	fmt.Println("\n************ gRPC query checking ************")
	nf := cq.GetNodeRes(conn)
	fmt.Println(nf.GetNodeInfo())
	//vf := cq.GetValRes(conn, viper.GetString("nova.val_addr"))
	//fmt.Println(vf.GetValidator().Tokens)
	wg.Wait()
}
