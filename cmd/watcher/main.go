package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	cfg "github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/rpc"
	rpctype "github.com/Carina-labs/HAL9000/rpc/types"
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const (
	NumWorker = 2
)

func main() {
	isTest := flag.Bool("test", false, "Decide whether it's test with localnet")
	cfg.SetChainInfo(*isTest)
	Nova := &cfg.NovaInfo{}
	Nova.Set("bot_addr")

	var wg sync.WaitGroup
	wg.Add(NumWorker)

	// Deprecated
	go func() {
		defer wg.Done()
		wsc, _, err := websocket.DefaultDialer.Dial(Nova.TmWsRPC.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		} else {
			log.Printf("connecting to %s", Nova.TmWsRPC.String())
		}
		defer func(c *websocket.Conn) {
			err := wsc.Close()
			if err != nil {
				utils.CheckErr(err, "", 1)
			}
		}(wsc)
		//myp := map[string]any{"query": "tm.event='Tx' And transfer.sender='nova1lds58drg8lvnaprcue2sqgfvjnz5ljlkq9lsyf'"}
		paramSet := map[string]any{"query": "tm.event='Tx'"}
		tmSubReq := &rpctype.RPCReq{JSONRPC: "2.0", Method: "subscribe", ID: "0", Params: paramSet}
		utils.CheckErr(err, "cannot marshal", 0)
		err = wsc.WriteJSON(tmSubReq)
		utils.CheckErr(err, "Cannot write JSON to Websocket : ", 0)

		for {
			var reply rpctype.RPCRes
			err = wsc.ReadJSON(&reply)
			evts := reply.Result.Events
			utils.CheckErr(err, "no reply from subscription", 1)
			//fmt.Println(evts)

			evt10, _ := rpc.CheckEvt(evts["nova.oracle.v1.ChainInfo.chain_id"])
			evt11, _ := rpc.CheckEvt(evts["nova.oracle.v1.ChainInfo.operator_address"])
			evt12, _ := rpc.CheckEvt(evts["nova.oracle.v1.ChainInfo.last_block_height"])
			evt13, _ := rpc.CheckEvt(evts["nova.oracle.v1.ChainInfo.app_hash"])

			oracleLog := fmt.Sprintf("Zone : %s, Operator : %s, Latest Block Height : %s Apphash : %s\n", evt10, evt11, evt12, evt13)
			totalLog := fmt.Sprintf("%s\n%s", time.Now().UTC().String(), oracleLog)
			_ = totalLog
			utils.CheckErr(err, "cannot write log to event.txt", 0)
		}
	}()

	// New implementation
	watchCtx, cancel := context.WithCancel(context.Background())
	go func() {
		defer wg.Done()

		novaWsc := rpc.MakeEventWatcher(*Nova.TmWsRPC)
		wsErr := novaWsc.Start()
		utils.CheckErr(wsErr, "", 0)
		query1 := fmt.Sprintf("tm.event='Tx' AND message.action='%s'", "/nova.oracle.v1.MsgUpdateChainState")
		//query1 := fmt.Sprintf("tm.event='Tx'")
		subsErr := novaWsc.Subscribe(watchCtx, query1)
		utils.CheckErr(subsErr, "", 0)

		parser := rpctype.NewTypedEventParser{}
		parser.SetProtoPkg("nova.oracle.v1")
		parser.SetProtoMsg("ChainInfo")

		for {
			res := <-novaWsc.ResponsesCh
			var wsRes rpctype.ResultEvent
			err := json.Unmarshal(res.Result, &wsRes)
			utils.CheckErr(err, "", 0)
			evts := wsRes.Events
			fmt.Printf("%s\n", evts[parser.EventWithFieldName("app_hash")])
			fmt.Printf("%s\n", evts[parser.EventWithFieldName("chain_id")])
			fmt.Printf("%s\n", evts[parser.EventWithFieldName("coin")])
			fmt.Printf("%s\n", evts[parser.EventWithFieldName("last_block_height")])
			fmt.Printf("%s\n", evts[parser.EventWithFieldName("decimal")])
			fmt.Printf("%s\n", evts[parser.EventWithFieldName("operator_address")])
		}

	}()

	wg.Wait()
	cancel()
}
