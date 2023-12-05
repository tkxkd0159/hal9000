package rpc

import (
	"net/url"

	tmrpc "github.com/tkxkd0159/HAL9000/rpc/types"
	"github.com/tkxkd0159/HAL9000/utils"
)

func MakeEventWatcher(remoteAddr url.URL) *tmrpc.WSClient {
	wsc, err := tmrpc.NewWS("//"+remoteAddr.Host, remoteAddr.Path)
	utils.CheckErr(err, "", 0)

	return wsc
}
