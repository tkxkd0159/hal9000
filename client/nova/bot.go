package nova

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	nova "github.com/Carina-labs/HAL9000/client/nova/types"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
)

var (
	wm sync.RWMutex
)

// GenTxByBot
// 1. Generate a TX with Msg (TxBuilder). If you set --generate-only, it makes unsigned tx and never broadcast
// 2. Sign the generated transaction with the keyring's account
// 3. Broadcast the tx to the Tendermint node using gPRC
func GenTxByBot(b *nova.Bot, onlyGen bool, msgs ...sdktypes.Msg) (ok bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Get panic while generating tx")
			ok = false
		}
	}()
	wm.Lock()
	defer wm.Unlock()
	if onlyGen {
		b.Ctx = b.Ctx.WithGenerateOnly(true)
	}

	err := tx.GenerateOrBroadcastTxWithFactory(b.Ctx, b.Txf, msgs...)
	if err != nil {
		if strings.Contains(err.Error(), "account sequence mismatch") {
			utils.LogErrWithFd(b.ErrLogger, err, "", 1)
			for {
				err = tx.GenerateOrBroadcastTxWithFactory(b.Ctx, b.Txf, msgs...)
				if !strings.Contains(err.Error(), "account sequence mismatch") {
					break
				}
				utils.LogErrWithFd(b.ErrLogger, err, "", ut.KEEP)
				time.Sleep(8 * time.Second)
			}
			if err != nil {
				utils.LogErrWithFd(b.ErrLogger, err, "something went wrong while make tx", 1)
				return false
			}
			return true

		}
		utils.LogErrWithFd(b.ErrLogger, err, "something went wrong while make tx", 1)
		return false
	}
	_, err = b.Ctx.Output.Write([]byte(fmt.Sprintf("%v: Tx was generated\n\n", time.Now())))
	utils.CheckErr(err, "cannot write log on output", 1)
	return true
}
