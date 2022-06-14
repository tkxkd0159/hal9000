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
	_, err := io.WriteString(w, "Hello user~")
	utils.CheckErr(err, "error occured while write get root response", ut.KEEP)
}

type apiHandler struct{}

func (apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/api/" {
		log.Println("got request /api/")
		utils.SetJSONHeader(w)
		resp := make(map[string]string)
		resp["message"] = "JSON test"
		resp["from"] = "bot"
		jsonResp, err := json.Marshal(resp)
		utils.CheckErr(err, "Error: can't JSON marshal", ut.KEEP)
		_, _ = w.Write(jsonResp)

	} else if req.URL.Path == "/api/1" {
		log.Println("got request /api/1")
		_, err := io.WriteString(w, "This is API 1 response")
		utils.CheckErr(err, "error occured while write response", ut.KEEP)
	} else if req.URL.Path == "/api/2" {
		log.Println("got request /api/2")
		_, err := io.WriteString(w, "This is API 2 response")
		utils.CheckErr(err, "error occured while write response", ut.KEEP)
	}
}
