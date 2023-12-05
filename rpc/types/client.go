package types

import (
	"context"
)

// Client describes the interface of Tendermint RPC client implementations.
type Client interface {
	// Start the client. Start must report an error if the client is running.
	Start() error

	// Stop the client. Stop must report an error if the client is not running.
	Stop() error

	// IsRunning reports whether the client is running.
	IsRunning() bool

	EventsClient
}

// EventsClient is reactive, you can subscribe to any message, given the proper
// string. see tendermint/types/events.go
type EventsClient interface {
	Subscribe(ctx context.Context, query string) error
	// Unsubscribe unsubscribes given subscriber from query.
	Unsubscribe(ctx context.Context, query string) error
	// UnsubscribeAll unsubscribes given subscriber from all the queries.
	UnsubscribeAll(ctx context.Context) error
}
