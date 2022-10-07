package logic

import (
	"log"
	"time"

	sdktype "github.com/cosmos/cosmos-sdk/types"

	"github.com/Carina-labs/HAL9000/client/nova"
	novam "github.com/Carina-labs/HAL9000/client/nova/msgs"
	novatypes "github.com/Carina-labs/HAL9000/client/nova/types"
	"github.com/Carina-labs/HAL9000/config"
)

func DepositGal(b *novatypes.Bot, host *config.HostChainInfo, denom string, amount int64) {

	i := 0
	intv := time.Duration(b.Interval)
	for {
		log.Printf("Pusher is ongoing for %d secs\n", int(intv)*i)

		msg1 := novam.MakeMsgDeposit(b.KrInfo.GetAddress(), b.KrInfo.GetAddress(), host.Name, denom, amount)
		msgs := []sdktype.Msg{msg1}
		nova.GenTxByBot(b, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}
