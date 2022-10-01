package msgs

import (
	"github.com/Carina-labs/nova/x/icacontrol/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func MakeMsgIcaAutoStaking(chainID string, operator sdktypes.AccAddress, decCoin sdktypes.DecCoin) *types.MsgIcaAutoStaking {
	coin, _ := sdktypes.NormalizeDecCoin(decCoin).TruncateDecimal()
	return types.NewMsgIcaAutoStaking(chainID, operator, coin)
}
