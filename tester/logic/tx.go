package logic

import (
	"github.com/Carina-labs/HAL9000/client/base"
	novam "github.com/Carina-labs/HAL9000/client/nova/msgs"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktype "github.com/cosmos/cosmos-sdk/types"
	"log"
	"os"
	"time"
)

func DepositGal(ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File, novaInfo IBCInfo, denom string, amount int64) {

	i := 0
	intv := time.Duration(interval)
	for {
		log.Printf("Pusher is ongoing for %d secs\n", int(intv)*i)

		msg1 := novam.MakeMsgDeposit(botInfo.GetAddress(), novaInfo.ZoneID, novaInfo.IBCPort, novaInfo.IBCChan, denom, amount)
		msgs := []sdktype.Msg{msg1}
		base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}
