package msgs

import (
	basev1beta "github.com/Carina-labs/nova/api/cosmos/base/v1beta1"
	"github.com/Carina-labs/nova/x/ibcstaking/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"log"
)

func MakeMsgIcaAutoStaking(chainID string, hostAddr string, operator sdktypes.AccAddress, decCoin *basev1beta.DecCoin) *types.MsgIcaAutoStaking {
	bigAmt, ok := sdktypes.NewIntFromString(decCoin.GetAmount())
	if !ok {
		log.Fatalln("Bigint conversion fail")
	}
	coin := sdktypes.NewCoin(decCoin.GetDenom(), bigAmt)
	return types.NewMsgIcaAutoStaking(chainID, hostAddr, operator, coin)
}
