package utils

import (
	"net/http"

	"google.golang.org/grpc"

	"github.com/tkxkd0159/HAL9000/utils/types"
)

func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func CloseGrpc(c *grpc.ClientConn) {
	err := c.Close()
	CheckErr(err, "", types.EXIT)
}
