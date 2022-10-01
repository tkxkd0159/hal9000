package msgs

import (
	"github.com/Carina-labs/nova/x/oracle/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"log"
)

func MakeMsgUpdateChainState(operator sdktypes.AccAddress, chainID string, denom string, amount string, blockHeight int64, apphash []byte) *types.MsgUpdateChainState {
	bigAmt, ok := sdktypes.NewIntFromString(amount)
	if !ok {
		log.Fatalln("Bigint conversion fail")
	}
	coin := sdktypes.NewCoin(denom, bigAmt)
	return types.NewMsgUpdateChainState(operator, chainID, coin, blockHeight, apphash)
}
