package nova

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/Carina-labs/HAL9000/client/base"
	nova "github.com/Carina-labs/HAL9000/client/nova/types"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
)

var (
	wm sync.RWMutex
)

func SetupBotKey(keyname, keyloc string, info *config.ChainNetInfo, bot config.BotScrt) {
	ctx := base.MakeContext(
		novaapp.ModuleBasics,
		bot.Address(),
		info.TmRPC.String(),
		info.ChainID,
		keyloc,
		keyring.BackendFile,
		os.Stdin,
		os.Stdout,
		false,
	)

	_ = base.MakeClientWithNewAcc(
		ctx,
		keyname,
		config.InputMnemonic(),
		sdktypes.FullFundraiserPath,
		hd.Secp256k1,
	)
}

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
				utils.LogErrWithFd(b.ErrLogger, err, "", 1)
				time.Sleep(6 * time.Second)
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
