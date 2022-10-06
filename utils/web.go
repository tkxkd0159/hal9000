package utils

import (
	"github.com/Carina-labs/HAL9000/utils/types"
	"google.golang.org/grpc"
	"net/http"
)

func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func CloseGrpc(c *grpc.ClientConn) {
	err := c.Close()
	CheckErr(err, "", types.EXIT)
}
