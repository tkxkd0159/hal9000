package api

import (
	"encoding/json"
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	"io"
	"log"
	"net/http"
)

func getRoot(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		return
	}
	log.Println("got request /")
	_, err := io.WriteString(w, "This is response")
	utils.CheckErr(err, "error occured while write get root response", ut.KEEP)
}

func getHello(w http.ResponseWriter, req *http.Request) {
	log.Println("got request /hello")
	utils.SetJSONHeader(w)
	resp := make(map[string]string)
	resp["message"] = "JSON: Hello user!"
	resp["from"] = "bot"
	jsonResp, err := json.Marshal(resp)
	utils.CheckErr(err, "Error: can't JSON marshal", ut.KEEP)
	_, _ = w.Write(jsonResp)
}

type apiHandler struct{}

func (apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/api/" {
		log.Println("got request /api/")
	} else if req.URL.Path == "/api/health" {
		_, err := io.WriteString(w, BotStatus.LastCommit.String())
		utils.CheckErr(err, "error occured while write response", ut.KEEP)
	}
}
