package logic

import (
	"fmt"
	"log"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/Carina-labs/HAL9000/config"
)

func botTickLog(botType string, elapsedSecs int, period int) {
	log.Printf("\n\n 🔁 %s Bot is ongoing for %d secs (run every %d secs) 🔁 \n\n", botType, elapsedSecs, period)
}

func botMsgLog(msgs []sdktypes.Msg) {
	for _, msg := range msgs {
		log.Printf(" 🔥 %s was sent 📦\n %s", sdktypes.MsgTypeURL(msg), msg.String())
	}
}

func initialBanner(botType string) {
	init := " _   _    _    _    ___   ___   ___   ___\n| | | |  / \\  | |  / _ \\ / _ \\ / _ \\ / _ \\\n| |_| | / _ \\ | | | (_) | | | | | | | | | |\n|  _  |/ ___ \\| |__\\__, | |_| | |_| | |_| |\n|_| |_/_/   \\_\\_____|/_/ \\___/ \\___/ \\___/\n"
	var after string
	switch botType {
	case config.ActOracle:
		after = "  ___  ____      _    ____ _     _____\n / _ \\|  _ \\    / \\  / ___| |   | ____|\n| | | | |_) |  / _ \\| |   | |   |  _|\n| |_| |  _ <  / ___ \\ |___| |___| |___\n \\___/|_| \\_\\/_/   \\_\\____|_____|_____|\n"
	case config.ActStake:
		after = " ____ _____  _    _  _______\n/ ___|_   _|/ \\  | |/ / ____|\n\\___ \\ | | / _ \\ | ' /|  _|\n ___) || |/ ___ \\| . \\| |___\n|____/ |_/_/   \\_\\_|\\_\\_____|\n"
	case config.ActAutoStake:
		after = "    _   _   _ _____ ___  ____ _____  _    _  _______\n   / \\ | | | |_   _/ _ \\/ ___|_   _|/ \\  | |/ / ____|\n  / _ \\| | | | | || | | \\___ \\ | | / _ \\ | ' /|  _|\n / ___ \\ |_| | | || |_| |___) || |/ ___ \\| . \\| |___\n/_/   \\_\\___/  |_| \\___/|____/ |_/_/   \\_\\_|\\_\\_____|\n"
	case config.ActWithdraw:
		after = "__        _____ _____ _   _ ____  ____      ___        __\n\\ \\      / /_ _|_   _| | | |  _ \\|  _ \\    / \\ \\      / /\n \\ \\ /\\ / / | |  | | | |_| | | | | |_) |  / _ \\ \\ /\\ / /\n  \\ V  V /  | |  | | |  _  | |_| |  _ <  / ___ \\ V  V /\n   \\_/\\_/  |___| |_| |_| |_|____/|_| \\_\\/_/   \\_\\_/\\_/\n"
	}
	fmt.Printf("%s%s\n\n", init, after)
}
