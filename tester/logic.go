package tester

import (
	"github.com/Carina-labs/HAL9000/client/base"
	galtype "github.com/Carina-labs/nova/x/gal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktype "github.com/cosmos/cosmos-sdk/types"
	"log"
	"os"
	"time"
)

func MakeMsgDeposit(from sdktype.AccAddress, zoneID, IBCPort, IBCChan, denom string, amount int64) *galtype.MsgDeposit {
	coin := sdktype.Coin{Denom: denom, Amount: sdktype.NewInt(amount)}
	return galtype.NewMsgDeposit(zoneID, from, coin, IBCPort, IBCChan)
}

func DepositGal(ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File, novaInfo IBCInfo, denom string, amount int64) {

	i := 0
	intv := time.Duration(interval)
	for {
		log.Printf("Pusher is ongoing for %d secs\n", int(intv)*i)

		msg1 := MakeMsgDeposit(botInfo.GetAddress(), novaInfo.ZoneID, novaInfo.IBCPort, novaInfo.IBCChan, denom, amount)
		msgs := []sdktype.Msg{msg1}
		base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}

type IBCInfo struct {
	ZoneID  string
	IBCPort string
	IBCChan string
}
