package ethhelpers

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// Deprecated:
	maxFilterLogWindow = 1000
)

type clientWithHTTPSubscriptions struct {
	Client
	createTicker func(context.Context, uint64) (BlockNumberTicker, error)
}

// TODO: Add options.

func NewClientWithHTTPSubscriptions(client Client, createTicker func(context.Context, uint64) (BlockNumberTicker, error)) Client {
	return &clientWithHTTPSubscriptions{
		Client:       client,
		createTicker: createTicker,
	}
}

// The context argument cancels the RPC request that sets up the subscription
// but has no effect on the subscription after Subscribe has returned.
func (c *clientWithHTTPSubscriptions) SubscribeFilterLogs(ctx context.Context, filterQuery ethereum.FilterQuery, logs chan<- types.Log) (ethereum.Subscription, error) {
	return SubscribeFilterLogsWithHTTP(ctx, &HTTPSubscriberOptions{
		Client:       c.Client,
		CreateTicker: c.createTicker,
		FilterQuery:  filterQuery,
		Logs:         logs,
	})
}
