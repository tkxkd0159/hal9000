package types

import (
	"errors"
	"github.com/tendermint/tendermint/libs/log"
	"sync"
)

var (
	// ErrClientRunning is returned by Start when the client is already running.
	ErrClientRunning = errors.New("client already running")

	// ErrClientNotRunning is returned by Stop when the client is not running.
	ErrClientNotRunning = errors.New("client is not running")
)

// RunState is a helper that a client implementation can embed to implement
// common plumbing for keeping track of run state and logging.
type RunState struct {
	Logger log.Logger

	mu        sync.Mutex
	name      string
	isRunning bool
	quit      chan struct{}
}

// NewRunState returns a new unstarted run state tracker with the given logging
// label and log sink. If logger == nil, a no-op logger is provided by default.
func NewRunState(name string, logger log.Logger) *RunState {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	return &RunState{
		name:   name,
		Logger: logger,
	}
}

// Start sets the state to running, or reports an error.
func (r *RunState) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.isRunning {
		r.Logger.Error("not starting client, it is already started", "client", r.name)
		return ErrClientRunning
	}
	r.Logger.Info("starting client", "client", r.name)
	r.isRunning = true
	r.quit = make(chan struct{})
	return nil
}

// Stop sets the state to not running, or reports an error.
func (r *RunState) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.isRunning {
		r.Logger.Error("not stopping client; it is already stopped", "client", r.name)
		return ErrClientNotRunning
	}
	r.Logger.Info("stopping client", "client", r.name)
	r.isRunning = false
	close(r.quit)
	return nil
}

// SetLogger updates the log sink.
func (r *RunState) SetLogger(logger log.Logger) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Logger = logger
}

// IsRunning reports whether the state is running.
func (r *RunState) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.isRunning
}

// Quit returns a channel that is closed when a call to Stop succeeds.
func (r *RunState) Quit() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.quit
}
