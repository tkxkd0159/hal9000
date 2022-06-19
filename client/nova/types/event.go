package types

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

type RpcReq struct {
	Jsonrpc string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
	ID      string         `json:"id"`
}

type RpcRes struct {
	ID      string    `json:"id"`
	JsonRPC string    `json:"jsonrpc"`
	Result  ResultRes `json:"result"`
}

type ResultRes struct {
	Query  string              `json:"query"`
	Events map[string][]string `json:"events"`
	Data   map[string]any      `json:"data"`
}

type Event = sdktypes.Event
type EvtAttr = []abcitypes.EventAttribute
