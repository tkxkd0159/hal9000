package types

import (
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type HostTransferChanID = string

type Bot struct {
	Ctx       client.Context
	Txf       tx.Factory
	KrInfo    keyring.Info
	Interval  int
	ErrLogger *os.File
	APIch     chan time.Time
}

func NewBot(ctx client.Context, txf tx.Factory, kr keyring.Info, interval int, errLogger *os.File, botch chan time.Time) *Bot {
	return &Bot{ctx, txf, kr, interval, errLogger, botch}
}
