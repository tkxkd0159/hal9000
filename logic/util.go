package logic

import (
	"log"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func botTickLog(botType string, elapsedSecs int) {
	log.Printf("----> %s Bot is ongoing for %d secs <----\n\n", botType, elapsedSecs)
}

func botMsgLog(msgs []sdktypes.Msg) {
	for _, msg := range msgs {
		log.Printf(" 🔥%s was sent 📦\n %s", sdktypes.MsgTypeURL(msg), msg.String())
	}
}
