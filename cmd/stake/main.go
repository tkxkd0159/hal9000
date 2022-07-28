package main

import (
	"flag"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/logic"
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
	sViper = cfg.Sviper
	common.SetBechPrefix()
}

func main() {
	isTest := flag.Bool("test", false, "Decide whether it's test with localnet")
	apiAddr := flag.String("api", "127.0.0.1:3336", "Set bot api address")
	keyname := flag.String("name", "nova_bot", "Set unique key name (uid)")
	newacc := flag.Bool("add", false, "Start client with making new account")
	chanID := flag.String("ch", "channel-0", "Nova Transfer Channel ID")
	hostchain := flag.String("host", "gaia", "Name of host chain registered with Nova")
	intv := flag.Int("interval", 10*60, "ibc-staking update interval (sec)")
	disp := flag.Bool("display", false, "Show context log through stdout")
	flag.Parse()
	flags := cfg.FlagOpts{Test: *isTest, New: *newacc, Disp: *disp, ExtIP: *apiAddr, Kn: *keyname, Host: *hostchain, Period: *intv, IBCChan: cfg.IBCChan{Nova: cfg.IBCPort{Transfer: *chanID}}}

	wg.Add(3)
	go func() {
		defer wg.Done()
		api.Server{}.On(flags.ExtIP)
	}()

	cfg.SetChainInfo(*isTest)
	krDir, logDir := cfg.SetInitialDir(*keyname, "logs/stake")
	fdLog, fdErr, fdErrExt := cfg.SetAllLogger(logDir, "ctxlog.txt", "nova_err.txt", "other_err.txt", flags.Disp)
	projFps := []*os.File{fdLog, fdErr, fdErrExt}
	defer func(fps ...*os.File) {
		for _, fp := range fps {
			err := fp.Close()
			utils.CheckErr(err, "", 1)
		}
	}(projFps...)

	// set pipe to ignore stdin tty
	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)
	novaBotAddr := viper.GetString("nova.bot_addr")
	novaIP := viper.GetString("net.ip.nova")
	novaTmAddr := novaIP + ":" + viper.GetString("net.port.tmrpc")
	novaTCPTmAddr := url.URL{Scheme: "tcp", Host: novaTmAddr}

	if flags.New {
		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			novaBotAddr,
			novaTCPTmAddr.String(),
			viper.GetString("nova.chain_id"),
			krDir,
			keyring.BackendFile,
			os.Stdin,
			fdLog,
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
		pp := cfg.GetPassphrase(sViper)
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
			fdLog,
			false,
		)
		os.Stdin = rpipe
		botInfo = common.LoadClientPubInfo(ctx, *keyname)
	}
	ctx = common.AddMoreFromInfo(ctx)
	txf := common.MakeTxFactory(ctx, "auto", "0unova", "", 1.1)

	go func(interval int) {
		defer wg.Done()
		logic.IcaStake(flags.Host, ctx, txf, botInfo, flags.IBCChan.Nova.Transfer, interval, fdErr)
	}(flags.Period)

	wg.Wait()

}
