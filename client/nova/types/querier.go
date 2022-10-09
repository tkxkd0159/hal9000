package types

import (
	galtypes "github.com/Carina-labs/nova/x/gal/types"
	icatypes "github.com/Carina-labs/nova/x/icacontrol/types"
)

type NovaQuerier interface {
	galQuerier
	icaQuerier
}

type galQuerier interface {
	CurrentDelegateVersion(zoneid string) *galtypes.QueryCurrentDelegateVersionResponse
	CurrentUndelegateVersion(zoneid string) *galtypes.QueryCurrentUndelegateVersionResponse
	CurrentWithdrawVersion(zoneid string) *galtypes.QueryCurrentWithdrawVersionResponse
}

type icaQuerier interface {
	CurrentAutoStakingVersion(zoneid string) *icatypes.QueryCurrentAutoStakingVersionResponse
}
