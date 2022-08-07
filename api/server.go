package api

import (
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	"net/http"
	"sync"
	"time"
)

type Server struct {
}

type TxCheckPoint struct {
	LastCommit time.Time
	mu         sync.Mutex
}

func (s *TxCheckPoint) SetCommitTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastCommit = t
}

var BotStatus *TxCheckPoint

func init() {
	BotStatus = &TxCheckPoint{LastCommit: time.Now()}
}

func (s Server) On(addr string) {
	http.Handle("/api/", apiHandler{})
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	err := http.ListenAndServe(addr, nil)
	utils.CheckErr(err, "cannot open http server", ut.EXIT)
}
