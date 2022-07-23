package main

import (
	"github.com/Carina-labs/HAL9000/client/common"
	novaTx "github.com/Carina-labs/HAL9000/client/nova/msgs"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"log"
	"os"
	"time"
)

func IcaStake(host string, txf tx.Factory, chanID string, interval int, errLogger *os.File) {

	stream := ut.Fstream{Err: errLogger}
	i := 0
	intv := time.Duration(interval)
	for {
		log.Printf("ICA-staking Bot is ongoing for %d secs\n", int(intv)*i)

		msg1 := novaTx.MakeMsgDelegate(host, botInfo.GetAddress(), "transfer", chanID)
		msgs := []sdktypes.Msg{msg1}
		common.GenTxWithFactory(stream, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}

}
