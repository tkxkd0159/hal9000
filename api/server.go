package api

import (
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	"net/http"
)

type Server struct{}

func (Server) On(addr string) {
	http.Handle("/api/", apiHandler{})
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	err := http.ListenAndServe(addr, nil)
	utils.CheckErr(err, "cannot open http server", ut.EXIT)
}
