package main

import (
	"flag"
	"fmt"
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/HAL9000/client/nova"
	nt "github.com/Carina-labs/HAL9000/client/nova/types"
	"github.com/Carina-labs/HAL9000/cmd"
	"github.com/Carina-labs/HAL9000/config"
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
	"reflect"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

var (
	ctx           client.Context
	botInfo       keyring.Info
	krDir, logDir string
	sViper        *viper.Viper
)

func init() {
	sViper = config.Sviper
	common.SetBechPrefix()
	krDir, logDir = cmd.SetInitialDir("/bot", "logs/stake")
}

func main() {
	//apiAddr := flag.String("api", "127.0.0.1:3335", "Set bot api address")
	keyname := flag.String("name", "nova-bot", "Set unique key name (uid)")
	newacc := flag.Bool("add", false, "Start client with making new account")
	disp := flag.Bool("display", false, "Show context log through stdout")
	flag.Parse()

	wg.Add(2)

	novaIP := viper.GetString("net.ip.nova")
	rawNovaTmAddr := novaIP + ":" + viper.GetString("net.port.tmrpc")
	novaTmAddr := url.URL{Scheme: "ws", Host: rawNovaTmAddr}
	u := url.URL{Scheme: "ws", Host: rawNovaTmAddr, Path: "/websocket"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			utils.CheckErr(err, "", 1)
		}
	}(c)

	// #### Start bot logic ####
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
	} else {
		fpErr = os.Stderr
		fpErrNova = os.Stderr
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

	// set pipe to ignore stdin tty
	rpipe, wpipe, err := os.Pipe()
	utils.CheckErr(err, "", 0)

	if *newacc {

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			viper.GetString("nova.bot_addr"),
			novaTmAddr.String(),
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
		pp := cmd.GetPassphrase(sViper)
		_, err = wpipe.Write([]byte(pp))
		utils.CheckErr(err, "", 0)

		ctx = common.MakeContext(
			novaapp.ModuleBasics,
			viper.GetString("nova.bot_addr"),
			novaTmAddr.String(),
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
	_ = txf
	_ = botInfo

	//myp := map[string]any{"query": "tm.event='Tx' And transfer.sender='nova1lds58drg8lvnaprcue2sqgfvjnz5ljlkq9lsyf'"}
	myp := map[string]any{"query": "tm.event='Tx'"}
	tmSubReq := &nt.RpcReq{Jsonrpc: "2.0", Method: "subscribe", ID: "0", Params: myp}
	utils.CheckErr(err, "cannot marshal", 0)
	err = c.WriteJSON(tmSubReq)
	utils.CheckErr(err, "Cannot write JSON to Websocket : ", 0)

	go func() {
		defer wg.Done()

		fp, err := os.OpenFile(path.Join(logDir, "event.txt"), os.O_CREATE|os.O_WRONLY, 0644)
		utils.CheckErr(err, "", 0)

		for {
			var reply nt.RpcRes
			err := c.ReadJSON(&reply)
			evts := reply.Result.Events
			utils.CheckErr(err, "no reply from subscription", 1)
			time.Sleep(3 * time.Second)

			evt1, _ := nova.CheckEvt(evts["nova.oracle.v1.ChainInfo.operator_address"])
			evt2, _ := nova.CheckEvt(evts["nova.oracle.v1.ChainInfo.last_block_height"])
			evt3, _ := nova.CheckEvt(evts["nova.intertx.v1.RegisteredZone.zone_name"])
			evt4, ok := nova.CheckEvt(evts["nova.intertx.v1.RegisteredZone.ica_connection_info"])
			if ok {

			}
			evt5, ok := nova.CheckEvt(evts["nova.gal.v1.DepositRecord.amount"])
			if ok {

			}

			oracleLog := fmt.Sprintf("Operator : %s, Latest Block Height : %s\n", evt1, evt2)
			icaLog := fmt.Sprintf("Zone name : %s, Controller Address : %s\n", evt3, evt4)
			galLog := fmt.Sprintf("User deposit amount : %s\n", evt5)

			if len(evts["nova.intertx.v1.RegisteredZone.ica_connection_info"]) > 0 {
				fmt.Println(reflect.TypeOf(evts["nova.intertx.v1.RegisteredZone.ica_connection_info"][0]))

			}
			totalLog := fmt.Sprintf("%s%s%s", oracleLog, icaLog, galLog)
			fmt.Print(totalLog)
			_, err = fmt.Fprintf(fp, "%s\n", totalLog)
			utils.CheckErr(err, "cannot write log to event.txt", 0)
		}
	}()

	wg.Wait()

}
