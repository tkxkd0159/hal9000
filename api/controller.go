package api

import (
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	"io"
	"log"
	"net/http"
)

type chkHandler struct{}

func NewChkHandler() *chkHandler {
	return &chkHandler{}
}

func (ch *chkHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/check/" {
		log.Println("You can check commit time")
	} else if req.URL.Path == "/check/commit" {
		_, err := io.WriteString(w, BotStatus.LastCommit.String())
		utils.CheckErr(err, "error occured while write response", ut.KEEP)
	}
}
