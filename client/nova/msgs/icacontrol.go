package msgs

import (
	"github.com/Carina-labs/nova/x/icacontrol/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func MakeMsgIcaAutoStaking(chainID string, operator sdktypes.AccAddress, decCoin sdktypes.DecCoin, seq uint64) *types.MsgIcaAutoStaking {
	coin, _ := sdktypes.NormalizeDecCoin(decCoin).TruncateDecimal()
	return func(name string, controllerAddr sdktypes.AccAddress, amount sdktypes.Coin, v uint64) *types.MsgIcaAutoStaking {
		return &types.MsgIcaAutoStaking{
			ZoneId:            name,
			ControllerAddress: controllerAddr.String(),
			Amount:            amount,
			Version:           v,
		}
	}(chainID, operator, coin, seq)
}
