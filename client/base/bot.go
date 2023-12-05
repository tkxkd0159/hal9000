package base

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	basetypes "github.com/tkxkd0159/HAL9000/client/base/types"
	"github.com/tkxkd0159/HAL9000/utils"
	ut "github.com/tkxkd0159/HAL9000/utils/types"
)

type TxErr int

const (
	NONE TxErr = iota
	NEXT
	NORMAL
	SEQMISMATCH
	CRITICAL
	REPEAT
)

const (
	SeqRecoverDelay    = time.Second * 4
	NormalTxRetryDelay = time.Second
)

var wm sync.Mutex

// GenTxByBot
// 1. Generate a TX with Msg (TxBuilder). If you set --generate-only, it makes unsigned tx and never broadcast
// -> If the transaction creation fails due to sequence mismatch, the transaction is regenerated again after the set recovery time.
// 2. Sign the generated transaction with the keyring's account
// 3. Broadcast the tx to the Tendermint node using gPRC
func GenTxByBot(b *basetypes.Bot, msgs ...sdktypes.Msg) (e TxErr) {
	defer func() {
		if err := recover(); err != nil {
			if realerr, ok := err.(error); ok {
				if errors.Is(realerr, ErrMustPanic) {
					panic(time.Now())
				}
			}

			log.Println(" ☠️ Get panic while generating tx ->", err)
			e = CRITICAL
		}
	}()

	wm.Lock()
	defer wm.Unlock()

	var err error
	var txbytes []byte
TXLOOP:
	for {
		txbytes, err = GenerateTx(b.Ctx, b.Txf, msgs...)
		status := handleTxErr(b.ErrLogger, err, b.Type)
		switch status {
		case NONE:
			break TXLOOP
		case NEXT:
			return NEXT
		case NORMAL:
			time.Sleep(NormalTxRetryDelay)
		case SEQMISMATCH:
			time.Sleep(SeqRecoverDelay)
		case REPEAT:
			return REPEAT
		}
	}

	for {
		err = BroadcastTx(b.Ctx, txbytes)
		if err == nil {
			break
		}

		utils.CheckErr(err, " ❌ something went wrong while broadcast tx", ut.KEEP)
	}

	_, err = b.Ctx.Output.Write([]byte(fmt.Sprintf("%v: Tx was generated\n\n", time.Now())))
	utils.CheckErr(err, " ❌ cannot write log on output", ut.KEEP)
	return NONE
}
