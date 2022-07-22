package msgs

import (
	"github.com/Carina-labs/HAL9000/client/common"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func MakeMsgSend(from any, to any, denoms []string, amounts []int64) (*banktypes.MsgSend, error) {

	from, err := common.CheckAccAddr(from)
	if err != nil {
		return nil, err
	}
	to, err = common.CheckAccAddr(to)
	if err != nil {
		return nil, err
	}

	var coins []sdktypes.Coin
	for i, denom := range denoms {
		c := sdktypes.NewCoin(denom, sdktypes.NewInt(amounts[i]))
		coins = append(coins, c)
	}

	return banktypes.NewMsgSend(from.(sdktypes.AccAddress), to.(sdktypes.AccAddress), sdktypes.NewCoins(coins...)), nil
}
