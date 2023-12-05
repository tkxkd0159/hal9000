package api

import (
	"io"
	"log"
	"net/http"

	"github.com/tkxkd0159/HAL9000/utils"
	ut "github.com/tkxkd0159/HAL9000/utils/types"
)

type ChkHandler struct{}

func NewChkHandler() *ChkHandler {
	return &ChkHandler{}
}

func (ch *ChkHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/check/" {
		log.Println("You can check commit time")
	} else if req.URL.Path == "/check/commit" {
		_, err := io.WriteString(w, BotStatus.LastCommit.String())
		utils.CheckErr(err, "error occurred while write response", ut.KEEP)
	}
}
