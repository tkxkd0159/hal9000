package msgs

import (
	"github.com/Carina-labs/nova/x/oracle/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"log"
)

func MakeMsgUpdateChainState(operator sdktypes.AccAddress, chainID string, amount string, denom string, decimal uint32, blockHeight int64, apphash []byte) *types.MsgUpdateChainState {
	bigAmt, ok := sdktypes.NewIntFromString(amount)
	if !ok {
		log.Fatalln("Bigint conversion fail")
	}
	coin := sdktypes.NewCoin(denom, bigAmt)
	return types.NewMsgUpdateChainState(operator, chainID, coin, decimal, blockHeight, apphash)
}
