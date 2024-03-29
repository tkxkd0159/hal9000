package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	cfg "github.com/tkxkd0159/HAL9000/config"
	"github.com/tkxkd0159/HAL9000/utils"
	ut "github.com/tkxkd0159/HAL9000/utils/types"
)

type Server struct{}

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

func Tracker(ch <-chan bool) {
	g := promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "current_commit_time",
			Help: "The last time the bot transaction was executed",
		},
	)
	for range ch {
		g.SetToCurrentTime()
		time.Sleep(time.Second * 2)
	}
}

func (s Server) On(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/check/", NewChkHandler())
	srv := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: time.Second * 3,
	}
	err := srv.ListenAndServe()
	utils.CheckErr(err, "cannot open http server", ut.EXIT)
}

func OpenMonitoringSrv(wg *sync.WaitGroup, datach <-chan time.Time, flags cfg.Observable) {
	defer wg.Done()
	relay := make(chan bool)
	go func() {
		for t := range datach {
			BotStatus.SetCommitTime(t)
			relay <- true
		}
	}()
	go Tracker(relay)
	Server{}.On(flags.GetExtIP())
}
