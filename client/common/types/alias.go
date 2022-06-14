package types

import (
	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	stakingv1beta1 "github.com/Carina-labs/nova/api/cosmos/staking/v1beta1"
)

// ###################  Base #####################

type NodeInfoReq = tendermintv1beta1.GetNodeInfoRequest
type NodeInfoRes = tendermintv1beta1.GetNodeInfoResponse

// ###################  Staking #####################

type ValInfoReq = stakingv1beta1.QueryValidatorRequest
type ValInfoRes = stakingv1beta1.QueryValidatorResponse
