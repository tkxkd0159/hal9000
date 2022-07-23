package main

import (
	"flag"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/cmd"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"sync"
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

func main() {
	isTest := flag.Bool("test", false, "Decide whether it's test with localnet")
	apiAddr := flag.String("api", "127.0.0.1:3337", "Set bot api address")
	keyname := flag.String("name", "nova_bot", "Set unique key name (uid)")
	newacc := flag.Bool("add", false, "Start client with making new account")
	chanID := flag.String("ch", "channel-41", "Host Transfer Channel ID")
	host := flag.String("host", "gaia", "Name of the host chain from which to obtain oracle info")
	intv := flag.Int("interval", 21*24*60*60, "Withdraw interval (sec)")
	disp := flag.Bool("display", false, "Show context log through stdout")
	flag.Parse()
	config.SetChainInfo(*isTest)

	novaBotAddr := viper.GetString("nova.bot_addr")
	novaIP := viper.GetString("net.ip.nova")
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
			utils.CheckErr(err, "", 1)

		}
	}(projFps...)

	// set pipe to ignore stdin tty
	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)

	if *newacc {

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			novaBotAddr,
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
			novaBotAddr,
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

	// ###### Start target bot logic ######
	go func(interval int) {
		defer wg.Done()
		UndelegateAndWithdraw(*host, txf, *chanID, interval, fpErrNova)
	}(*intv)

	wg.Wait()
}
