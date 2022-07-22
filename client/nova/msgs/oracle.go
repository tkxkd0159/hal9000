package msgs

import (
	"github.com/Carina-labs/HAL9000/client/common"
	"github.com/Carina-labs/nova/x/oracle/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func MakeMsgUpdateChainState(operator any, denom string, amount int64, decimal uint64, bheight uint64) (*types.MsgUpdateChainState, error) {
	operator, err := common.CheckAccAddr(operator)
	if err != nil {
		return nil, err
	}
	coin := sdktypes.NewCoin(denom, sdktypes.NewInt(amount))
	return types.NewMsgUpdateChainState(coin, operator.(sdktypes.AccAddress), decimal, bheight), nil
}
