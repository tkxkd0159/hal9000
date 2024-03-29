package msgs

import (
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/Carina-labs/nova/x/gal/types"
)

func MakeMsgDelegate(chainID string, operator sdktypes.AccAddress, seq uint64, ibctimeout uint64) *types.MsgDelegate {
	return types.NewMsgDelegate(chainID, seq, operator, ibctimeout)
}

func MakeMsgIcaWithdraw(chainID string, operator sdktypes.AccAddress, portID string, chanID string, blockTS time.Time, seq uint64, ibctimeout uint64) *types.MsgIcaWithdraw {
	return types.NewMsgIcaWithdraw(chainID, operator, portID, chanID, blockTS, seq, ibctimeout)
}

func MakeMsgUndelegate(chainID string, operator sdktypes.AccAddress, seq uint64, ibctimeout uint64) *types.MsgUndelegate {
	return types.NewMsgUndelegate(chainID, seq, operator, ibctimeout)
}

func MakeMsgDeposit(from, claimer sdktypes.AccAddress, zoneID, denom string, amount int64, ibctimeout uint64) *types.MsgDeposit {
	coin := sdktypes.Coin{Denom: denom, Amount: sdktypes.NewInt(amount)}
	return types.NewMsgDeposit(zoneID, from, claimer, coin, ibctimeout)
}

func MakeMsgClaimAllSnAsset(zoneID string, from sdktypes.AccAddress) *types.MsgClaimAllSnAsset {
	return types.NewMsgClaimAllSnAsset(zoneID, from)
}
