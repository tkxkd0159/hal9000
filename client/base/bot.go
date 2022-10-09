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
	NORMAL
	SEQMISMATCH
)

var (
	wm sync.Mutex
)

// GenTxByBot
// 1. Generate a TX with Msg (TxBuilder). If you set --generate-only, it makes unsigned tx and never broadcast
// 2. Sign the generated transaction with the keyring's account
// 3. Broadcast the tx to the Tendermint node using gPRC
func GenTxByBot(b *basetypes.Bot, onlyGen bool, msgs ...sdktypes.Msg) (ok bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Get panic while generating tx ->", err)
			ok = false
		}
	}()

	wm.Lock()
	defer wm.Unlock()
	if onlyGen {
		b.Ctx = b.Ctx.WithGenerateOnly(true)
	}
	var err error
	var txbytes []byte
TXLOOP:
	for {
		txbytes, err = GenerateTx(b.Ctx, b.Txf, msgs...)
		status := handleSeqErr(b.ErrLogger, err)
		switch status {
		case NONE:
			break TXLOOP
		case NORMAL:
		case SEQMISMATCH:
			time.Sleep(time.Second * 4)
		}
	}

	for {
		err = BroadcastTx(b.Ctx, txbytes)
		if err == nil {
			break
		} else {
			utils.CheckErr(err, "something went wrong while broadcast tx", ut.KEEP)
		}
	}

	_, err = b.Ctx.Output.Write([]byte(fmt.Sprintf("%v: Tx was generated\n\n", time.Now())))
	utils.CheckErr(err, "cannot write log on output", ut.KEEP)
	return true
}

func handleSeqErr(f *os.File, e error) TxErr {
	if e != nil {
		if strings.Contains(e.Error(), "account sequence mismatch") {
			utils.LogErrWithFd(f, e, " ❌ ", ut.KEEP)
			return SEQMISMATCH
		}
		utils.LogErrWithFd(f, e, " ❌ something went wrong while generate tx", ut.KEEP)
		return NORMAL
	}
	return NONE
}
