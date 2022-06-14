package common

import (
	"errors"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type (
	AccAddr = sdktypes.AccAddress
)

func CheckAccAddr(target any) (AccAddr, error) {
	switch target := target.(type) {
	case AccAddr:
		return target, nil
	case string:
		addr, err := sdktypes.AccAddressFromBech32(target)
		if err != nil {
			return nil, err
		}
		return addr, nil
	case []byte:
		return target, nil
	default:
		return nil, errors.New("cannot covert target to AccAddress")
	}
}

// ###################  Bank #####################

func MakeSendMsg(from any, to any, denoms []string, amounts []int64) (*banktypes.MsgSend, error) {

	from, err := CheckAccAddr(from)
	if err != nil {
		return nil, err
	}
	to, err = CheckAccAddr(to)
	if err != nil {
		return nil, err
	}

	var coins []sdktypes.Coin
	for i, denom := range denoms {
		c := sdktypes.NewCoin(denom, sdktypes.NewInt(amounts[i]))
		coins = append(coins, c)
	}

	return banktypes.NewMsgSend(from.(AccAddr), to.(AccAddr), sdktypes.NewCoins(coins...)), nil
}
