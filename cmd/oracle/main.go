package main

import (
	"flag"
	"fmt"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	commonTx "github.com/Carina-labs/HAL9000/client/common/msgs"
	"github.com/Carina-labs/HAL9000/client/common/query"
	novaTx "github.com/Carina-labs/HAL9000/client/nova/msgs"
	"github.com/Carina-labs/HAL9000/cmd"
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
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

var (
	ctx     client.Context
	botInfo keyring.Info
	sViper  *viper.Viper
)

func init() {
	sViper = config.Sviper
	common.SetBechPrefix()
}

// FIXME: wasmvm doesn't support AArch64. Need to set GOARCH=amd64
func main() {
	isTest := flag.Bool("test", false, "Decide whether it's test with localnet")
	apiAddr := flag.String("api", "127.0.0.1:3334", "Set bot api address")
	keyname := flag.String("name", "nova_bot", "Set unique key name (uid)")
	newacc := flag.Bool("add", false, "Start client with making new account")
	intv := flag.Int("interval", 5, "Oracle update interval (sec)")
	disp := flag.Bool("display", false, "Show context log through stdout")
	flag.Parse()
	config.SetChainInfo(*isTest)

	novaIP := viper.GetString("net.ip.nova")
	novaGrpcAddr := novaIP + ":" + viper.GetString("net.port.grpc")
	novaTCPTmAddr := &url.URL{Scheme: "tcp", Host: novaIP + ":" + viper.GetString("net.port.tmrpc")}

	// Open api endpoint to check bot
	wg.Add(2)
	go func() {
		defer wg.Done()
		api.Server{}.On(*apiAddr)
	}()

	krDir, logDir := cmd.SetInitialDir(*keyname, "logs/oracle")
	fpLog, fpErr, fpErrNova := cmd.SetAllLogger(logDir, "ctxlog.txt", "nova_err.txt", "other_err.txt", disp)

	projFps := []*os.File{fpLog, fpErr, fpErrNova}
	defer func(fps ...*os.File) {
		for _, fp := range fps {
			err := fp.Close()
			if err != nil {
				utils.CheckErr(err, "", 1)
			}
		}
	}(projFps...)

	// set pipe to ignore stdin tty
	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)

	// #### Start bot logic ####

	if *newacc {

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			viper.GetString("nova.bot_addr"),
			novaTCPTmAddr.String(),
			viper.GetString("nova.chain_id"),
			krDir,
			keyring.BackendFile,
			os.Stdin,
			fpLog,
			false,
		)

		botInfo = common.MakeClientWithNewAcc(
			ctx,
			*keyname,
			sViper.GetString(*keyname),
			sdktypes.FullFundraiserPath,
			hd.Secp256k1,
		)
		os.Exit(0)
	} else {
		pp := cmd.GetPassphrase(sViper)
		_, err = wpipe.Write([]byte(pp))
		utils.CheckErr(err, "", 0)

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			viper.GetString("nova.bot_addr"),
			novaTCPTmAddr.String(),
			viper.GetString("nova.chain_id"),
			krDir,
			keyring.BackendFile,
			rpipe,
			fpLog,
			false,
		)
		os.Stdin = rpipe

		botInfo = common.LoadClientPubInfo(ctx, *keyname)
	}
	ctx = common.AddMoreFromInfo(ctx)
	txf := common.MakeTxFactory(ctx, "auto", "0unova", "", 1.1)

	// ** Build TX
	go func(interval int) {
		defer wg.Done()

		stream := ut.Fstream{Err: fpErrNova}
		i := 0
		intv := time.Duration(interval)
		for {
			log.Printf("Bot is ongoing for %d secs\n", int(intv)*i)

			msg1, err := commonTx.MakeMsgSend(botInfo.GetAddress(), viper.GetString("nova.target_addr"), []string{"unova"}, []int64{777})
			if err != nil {
				utils.CheckErr(err, "", 0)
				continue
			}
			msg2, err := novaTx.MakeMsgUpdateChainState(botInfo.GetAddress(), "uatom", 7654321, 6, 777)
			if err != nil {
				utils.CheckErr(err, "", 0)
				continue
			}
			msgs := []sdktypes.Msg{msg1, msg2}

			common.GenTxWithFactory(stream, ctx, txf, false, msgs...)
			time.Sleep(intv * time.Second)
			i++
		}
	}(*intv)

	fmt.Println("\n************ gRPC query checking ************")
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

	cq := &query.CosmosQueryClient{ClientConn: conn}
	fmt.Println(cq.GetNodeRes().GetNodeInfo())
	nv := cq.GetValInfo(viper.GetString("nova.val_addr"))
	nb := cq.GetBlockByHeight(10)
	nbh := nb.GetBlock().Header

	st := fmt.Sprintf("Staked nova on Our validator : %s\n", nv.GetValidator().Tokens)
	proof := fmt.Sprintf("chain ID : %s, height : %d, apphash : %s, proposer : %s \n",
		nbh.ChainId, nbh.Height, utils.B64ToStr(nbh.AppHash), utils.B64ToStr(nbh.ProposerAddress))
	fmt.Printf("%s %s", st, proof)
	wg.Wait()
}
