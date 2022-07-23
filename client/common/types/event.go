package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type RPCReq struct {
	JSONRPC string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
	ID      string         `json:"id"`
}

type RPCRes struct {
	ID      string    `json:"id"`
	JSONRPC string    `json:"jsonrpc"`
	Result  ResultRes `json:"result"`
}

type ResultRes struct {
	Query  string              `json:"query"`
	Events map[string][]string `json:"events"`
	Data   map[string]any      `json:"data"`
}

type Event = sdktypes.Event
type EvtAttr = []abcitypes.EventAttribute
