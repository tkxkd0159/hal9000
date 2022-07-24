package main

import (
	"flag"
	"fmt"
	"github.com/Carina-labs/HAL9000/api"
	"github.com/Carina-labs/HAL9000/client/common"
	nt "github.com/Carina-labs/HAL9000/client/common/types"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
	"path"
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
	sViper = cfg.Sviper
	common.SetBechPrefix()
}

func main() {
	isTest := flag.Bool("test", false, "Decide whether it's test with localnet")
	apiAddr := flag.String("api", "127.0.0.1:3336", "Set bot api address")
	keyname := flag.String("name", "nova_bot", "Set unique key name (uid)")
	newacc := flag.Bool("add", false, "Start client with making new account")
	chanID := flag.String("ch", "channel-0", "Nova Transfer Channel ID")
	hostchain := flag.String("host", "gaia", "Name of the host chain from which to obtain oracle info")
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
	novaWsTmAddr := url.URL{Scheme: "ws", Host: novaTmAddr, Path: "/websocket"}

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

	wsc, _, err := websocket.DefaultDialer.Dial(novaWsTmAddr.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	} else {
		log.Printf("connecting to %s", novaWsTmAddr.String())
	}
	defer func(c *websocket.Conn) {
		err := wsc.Close()
		if err != nil {
			utils.CheckErr(err, "", 1)
		}
	}(wsc)

	//myp := map[string]any{"query": "tm.event='Tx' And transfer.sender='nova1lds58drg8lvnaprcue2sqgfvjnz5ljlkq9lsyf'"}
	myp := map[string]any{"query": "tm.event='Tx'"}
	tmSubReq := &nt.RPCReq{JSONRPC: "2.0", Method: "subscribe", ID: "0", Params: myp}
	utils.CheckErr(err, "cannot marshal", 0)
	err = wsc.WriteJSON(tmSubReq)
	utils.CheckErr(err, "Cannot write JSON to Websocket : ", 0)

	// ###### Start target bot logic ######
	go func() {
		defer wg.Done()

		fp, err := os.OpenFile(path.Join(logDir, "event.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		utils.CheckErr(err, "", 0)

		for {
			var reply nt.RPCRes
			err = wsc.ReadJSON(&reply)
			evts := reply.Result.Events
			utils.CheckErr(err, "no reply from subscription", 1)

			evt10, _ := common.CheckEvt(evts["nova.oracle.v1.ChainInfo.chain_id"])
			evt11, _ := common.CheckEvt(evts["nova.oracle.v1.ChainInfo.operator_address"])
			evt12, _ := common.CheckEvt(evts["nova.oracle.v1.ChainInfo.last_block_height"])
			evt13, _ := common.CheckEvt(evts["nova.oracle.v1.ChainInfo.app_hash"])

			oracleLog := fmt.Sprintf("Zone : %s, Operator : %s, Latest Block Height : %s Apphash : %s\n", evt10, evt11, evt12, evt13)

			totalLog := fmt.Sprintf("%s\n%s", time.Now().UTC().String(), oracleLog)
			fmt.Print(totalLog)
			_, err = fmt.Fprintf(fp, "%s\n", totalLog)
			utils.CheckErr(err, "cannot write log to event.txt", 0)
		}
	}()

	go func(interval int) {
		defer wg.Done()
		IcaStake(flags.Host, txf, flags.IBCChan.Nova.Transfer, interval, fdErr)
	}(flags.Period)

	wg.Wait()

}
