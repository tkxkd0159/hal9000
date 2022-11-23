package base

import (
	"os"
	"strings"

	galtypes "github.com/Carina-labs/nova/x/gal/types"
	oracletypes "github.com/Carina-labs/nova/x/oracle/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
)

func handleTxErr(f *os.File, e error, bottype string) TxErr {
	if e != nil {
		if strings.Contains(e.Error(), sdkerr.ErrWrongSequence.Error()) {
			utils.LogErrWithFd(f, e, " ❌ ", ut.KEEP)
			return SEQMISMATCH
		} else if strings.Contains(e.Error(), galtypes.ErrInvalidIcaVersion.Error()) {
			utils.LogErrWithFd(f, e, " ❌ ica sequence was not updated yet when the bot queried\n", ut.KEEP)
			return REPEAT
		} else if strings.Contains(e.Error(), sdkerr.ErrInvalidAddress.Error()) {
			panic(sdkerr.Wrap(e, "please check your controller address in the keyring"))
		}

		switch bottype {
		case config.ActOracle:
			if strings.Contains(e.Error(), oracletypes.ErrInvalidBlockHeight.Error()) {
				utils.LogErrWithFd(f, e, " ❌ oracle info was outdated due to the oracle bot's update. It will regenerate tx\n", ut.KEEP)
				return REPEAT
			}
		case config.ActStake:
			if strings.Contains(e.Error(), galtypes.ErrNoDepositRecord.Error()) || strings.Contains(e.Error(), galtypes.ErrInsufficientFunds.Error()) {
				utils.LogErrWithFd(f, e, " ❌ There is no asset to delegate on this host zone  ➡️ go to next batch\n", ut.KEEP)
				return NEXT
			}
		case config.ActUndelegate, config.ActWithdraw:
			if strings.Contains(e.Error(), "no coins to undelegate") {
				utils.LogErrWithFd(f, e, " ❌ There is no asset to undelegate on this host zone  ➡️ go to next batch\n", ut.KEEP)
				return NEXT
			} else if strings.Contains(e.Error(), galtypes.ErrCanNotWithdrawAsset.Error()) {
				utils.LogErrWithFd(f, e, " ❌ There is no asset to withdraw on this host zone  ➡️ go to next batch\n", ut.KEEP)
				return NEXT
			}
		}

		utils.LogErrWithFd(f, e, " ❌ something went wrong while generate tx", ut.KEEP)
		return NORMAL
	}
	return NONE
}
