package types

import (
	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	txv1beta1 "github.com/cosmos/cosmos-sdk/types/tx"
)

type CommonQuerier interface {
	baseQuerier
	tmQuerier
	stakeQuerier
}

type baseQuerier interface {
	GetTx(hash string) *txv1beta1.GetTxResponse
}

type tmQuerier interface {
	GetNodeRes() *NodeInfoRes
	GetBlockByHeight(height int64) *tendermintv1beta1.GetBlockByHeightResponse
	GetLatestBlock() *tendermintv1beta1.GetLatestBlockResponse
}

type stakeQuerier interface {
	GetValInfo(valAddr string) *ValInfoRes
}
