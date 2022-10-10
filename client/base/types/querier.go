package types

import (
	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	stakingv1beta1 "github.com/Carina-labs/nova/api/cosmos/staking/v1beta1"
	bankv1beta1 "github.com/cosmos/cosmos-sdk/x/bank/types"
	distv1beta1 "github.com/cosmos/cosmos-sdk/x/distribution/types"

	txv1beta1 "github.com/cosmos/cosmos-sdk/types/tx"
)

type BaseQuerier interface {
	txQuerier
	tmQuerier
	bankQuerier
	stakeQuerier
	distQuerier
}

type bankQuerier interface {
	GetBalance(address string, denom string) (*bankv1beta1.QueryBalanceResponse, error)
}

type txQuerier interface {
	GetTx(hash string) (*txv1beta1.GetTxResponse, error)
}

type tmQuerier interface {
	GetNodeRes() (*tendermintv1beta1.GetNodeInfoResponse, error)
	GetBlockByHeight(height int64) (*tendermintv1beta1.GetBlockByHeightResponse, error)
	GetLatestBlock() (*tendermintv1beta1.GetLatestBlockResponse, error)
}

type stakeQuerier interface {
	GetValInfo(valAddr string) (*stakingv1beta1.QueryValidatorResponse, error)
	GetHistoricalInfo(height int64) (*stakingv1beta1.QueryHistoricalInfoResponse, error)
	GetDelegation(delAddr, valAddr string) (*stakingv1beta1.QueryDelegationResponse, error)
}

type distQuerier interface {
	GetRewards(delAddr string, valAddr string) (*distv1beta1.QueryDelegationRewardsResponse, error)
}
