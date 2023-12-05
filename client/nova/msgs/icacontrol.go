package msgs

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/Carina-labs/nova/x/icacontrol/types"
)

func MakeMsgIcaAutoStaking(chainID string, operator sdktypes.AccAddress, decCoin sdktypes.DecCoin, seq uint64, ibctimeout uint64) *types.MsgIcaAutoStaking {
	coin, _ := sdktypes.NormalizeDecCoin(decCoin).TruncateDecimal()
	return func(name string, controllerAddr sdktypes.AccAddress, amount sdktypes.Coin, v uint64, t uint64) *types.MsgIcaAutoStaking {
		return &types.MsgIcaAutoStaking{
			ZoneId:            name,
			ControllerAddress: controllerAddr.String(),
			Amount:            amount,
			Version:           v,
			TimeoutTimestamp:  t,
		}
	}(chainID, operator, coin, seq, ibctimeout)
}
