package types

import (
	"context"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

// Client describes the interface of Tendermint RPC client implementations.
type Client interface {
	// These methods define the operational structure of the client.

	// Start the client. Start must report an error if the client is running.
	Start() error

	// Stop the client. Stop must report an error if the client is not running.
	Stop() error

	// IsRunning reports whether the client is running.
	IsRunning() bool

	ABCIClient
	EventsClient
	HistoryClient
	SignClient
	StatusClient
	EvidenceClient
	MempoolClient
}

type ABCIQueryOptions struct {
	Height int64
	Prove  bool
}

// DefaultABCIQueryOptions are latest height (0) and prove false.
var DefaultABCIQueryOptions = ABCIQueryOptions{Height: 0, Prove: false}

// ABCIClient groups together the functionality that principally affects the
// ABCI app.
//
// In many cases this will be all we want, so we can accept an interface which
// is easier to mock.
type ABCIClient interface {
	// Reading from abci app

	ABCIInfo(context.Context) (*ResultABCIInfo, error)
	ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*ResultABCIQuery, error)
	ABCIQueryWithOptions(ctx context.Context, path string, data bytes.HexBytes,
		opts ABCIQueryOptions) (*ResultABCIQuery, error)

	// Writing to abci app

	BroadcastTxCommit(context.Context, types.Tx) (*ResultBroadcastTxCommit, error)
	BroadcastTxAsync(context.Context, types.Tx) (*ResultBroadcastTx, error)
	BroadcastTxSync(context.Context, types.Tx) (*ResultBroadcastTx, error)
}

// SignClient groups together the functionality needed to get valid signatures
// and prove anything about the chain.
type SignClient interface {
	Block(ctx context.Context, height *int64) (*ResultBlock, error)
	BlockByHash(ctx context.Context, hash bytes.HexBytes) (*ResultBlock, error)
	BlockResults(ctx context.Context, height *int64) (*ResultBlockResults, error)
	Header(ctx context.Context, height *int64) (*ResultHeader, error)
	HeaderByHash(ctx context.Context, hash bytes.HexBytes) (*ResultHeader, error)
	Commit(ctx context.Context, height *int64) (*ResultCommit, error)
	Validators(ctx context.Context, height *int64, page, perPage *int) (*ResultValidators, error)
	Tx(ctx context.Context, hash bytes.HexBytes, prove bool) (*ResultTx, error)

	// TxSearch defines a method to search for a paginated set of transactions by
	// DeliverTx event search criteria.
	TxSearch(
		ctx context.Context,
		query string,
		prove bool,
		page, perPage *int,
		orderBy string,
	) (*ResultTxSearch, error)

	// BlockSearch defines a method to search for a paginated set of blocks by
	// BeginBlock and EndBlock event search criteria.
	BlockSearch(
		ctx context.Context,
		query string,
		page, perPage *int,
		orderBy string,
	) (*ResultBlockSearch, error)
}

// HistoryClient provides access to data from genesis to now in large chunks.
type HistoryClient interface {
	Genesis(context.Context) (*ResultGenesis, error)
	GenesisChunked(context.Context, uint) (*ResultGenesisChunk, error)
	BlockchainInfo(ctx context.Context, minHeight, maxHeight int64) (*ResultBlockchainInfo, error)
}

// StatusClient provides access to general chain info.
type StatusClient interface {
	Status(context.Context) (*ResultStatus, error)
}

// EventsClient is reactive, you can subscribe to any message, given the proper
// string. see tendermint/types/events.go
type EventsClient interface {
	// Subscribe subscribes given subscriber to query. Returns a channel with
	// cap=1 onto which events are published. An error is returned if it fails to
	// subscribe. outCapacity can be used optionally to set capacity for the
	// channel. Channel is never closed to prevent accidental reads.
	//
	// ctx cannot be used to unsubscribe. To unsubscribe, use either Unsubscribe
	// or UnsubscribeAll.
	Subscribe(ctx context.Context, subscriber, query string, outCapacity ...int) (out <-chan ResultEvent, err error)
	// Unsubscribe unsubscribes given subscriber from query.
	Unsubscribe(ctx context.Context, subscriber, query string) error
	// UnsubscribeAll unsubscribes given subscriber from all the queries.
	UnsubscribeAll(ctx context.Context, subscriber string) error
}

// MempoolClient shows us data about current mempool state.
type MempoolClient interface {
	UnconfirmedTxs(ctx context.Context, limit *int) (*ResultUnconfirmedTxs, error)
	NumUnconfirmedTxs(context.Context) (*ResultUnconfirmedTxs, error)
	CheckTx(context.Context, types.Tx) (*ResultCheckTx, error)
	RemoveTx(context.Context, types.TxKey) error
}

// EvidenceClient is used for submitting an evidence of the malicious
// behavior.
type EvidenceClient interface {
	BroadcastEvidence(context.Context, types.Evidence) (*ResultBroadcastEvidence, error)
}

// RemoteClient is a Client, which can also return the remote network address.
type RemoteClient interface {
	Client

	// Remote returns the remote network address in a string form.
	Remote() string
}
