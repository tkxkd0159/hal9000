package msgs

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/Carina-labs/HAL9000/client/base"
)

func MakeMsgSend(from any, to any, denoms []string, amounts []int64) (*banktypes.MsgSend, error) {

	from, err := base.CheckAccAddr(from)
	if err != nil {
		return nil, err
	}
	to, err = base.CheckAccAddr(to)
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
