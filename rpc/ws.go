package rpc

import (
	"net/url"

	tmrpc "github.com/Carina-labs/HAL9000/rpc/types"
	"github.com/Carina-labs/HAL9000/utils"
)

func MakeEventWatcher(remoteAddr url.URL) *tmrpc.WSClient {
	wsc, err := tmrpc.NewWS("//"+remoteAddr.Host, remoteAddr.Path)
	utils.CheckErr(err, "", 0)
	//wsc.SetLogger(tmlog.NewTMJSONLogger(tmlog.NewSyncWriter(os.Stdout)))
	return wsc
}
