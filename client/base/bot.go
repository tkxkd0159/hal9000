package base

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	basetypes "github.com/Carina-labs/HAL9000/client/base/types"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
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

var (
	wm sync.Mutex
)

// GenTxByBot
// 1. Generate a TX with Msg (TxBuilder). If you set --generate-only, it makes unsigned tx and never broadcast
// -> If the transaction creation fails due to sequence mismatch, the transaction is regenerated again after the set recovery time.
// 2. Sign the generated transaction with the keyring's account
// 3. Broadcast the tx to the Tendermint node using gPRC
func GenTxByBot(b *basetypes.Bot, msgs ...sdktypes.Msg) (e TxErr) {
	defer func() {
		if err := recover(); err != nil {
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
		status := handleTxErr(b.ErrLogger, err)
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
		} else {
			utils.CheckErr(err, " ❌ something went wrong while broadcast tx", ut.KEEP)
		}
	}

	_, err = b.Ctx.Output.Write([]byte(fmt.Sprintf("%v: Tx was generated\n\n", time.Now())))
	utils.CheckErr(err, " ❌ cannot write log on output", ut.KEEP)
	return NONE
}

func handleTxErr(f *os.File, e error) TxErr {
	if e != nil {
		if strings.Contains(e.Error(), "account sequence mismatch") {
			utils.LogErrWithFd(f, e, " ❌ ", ut.KEEP)
			return SEQMISMATCH
		} else if strings.Contains(e.Error(), "cannot change state") {
			utils.LogErrWithFd(f, e, " ❌ There is no asset to delegate on this host zone  ➡️ go to next batch\n", ut.KEEP)
			return NEXT
		} else if strings.Contains(e.Error(), "invalid coins") {
			utils.LogErrWithFd(f, e, " ❌ There is no reward to autostake on this host zone  ➡️ go to next batch\n", ut.KEEP)
			return NEXT
		} else if strings.Contains(e.Error(), "no coins to undelegate") {
			utils.LogErrWithFd(f, e, " ❌ There is no asset to undelegate on this host zone  ➡️ go to next batch\n", ut.KEEP)
			return NEXT
		} else if strings.Contains(e.Error(), "cannot withdraw funds") {
			utils.LogErrWithFd(f, e, " ❌ There is no asset to withdraw on this host zone  ➡️ go to next batch\n", ut.KEEP)
			return NEXT
		} else if strings.Contains(e.Error(), "current block height must be higher than the previous block height") {
			utils.LogErrWithFd(f, e, " ❌ oracle info was outdated due to the oracle bot's update. It will regenerate tx\n", ut.KEEP)
			return REPEAT
		} else if strings.Contains(e.Error(), "invalid ica version") {
			utils.LogErrWithFd(f, e, " ❌ ica sequence was not updated yet when the bot queried\n", ut.KEEP)
			return REPEAT
		}

		utils.LogErrWithFd(f, e, " ❌ something went wrong while generate tx", ut.KEEP)
		return NORMAL
	}
	return NONE
}
