package ethhelpers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	maxFilterLogWindow = 1000
)

type clientWithHTTPSubscriptions struct {
	Client
	newTickerFn func(context.Context, BlockNumberReader) BlockNumberTicker
}

// TODO: Add options.

func NewClientWithHTTPSubscriptions(client Client, newTickerFn func(context.Context, BlockNumberReader) BlockNumberTicker) Client {
	return &clientWithHTTPSubscriptions{
		Client:      client,
		newTickerFn: newTickerFn,
	}
}

// The context argument cancels the RPC request that sets up the subscription
// but has no effect on the subscription after Subscribe has returned.
func (c *clientWithHTTPSubscriptions) SubscribeFilterLogs(ctx context.Context, filterQuery ethereum.FilterQuery, logChan chan<- types.Log) (ethereum.Subscription, error) {
	if logChan == nil {
		panic("channel given to SubscribeFilterLogs must not be nil")
	}

	subCtx, subCancel := context.WithCancel(context.Background())
	errChan := make(chan error, 1)

	sub := &clientSubscription{
		subCancel: subCancel,
		errChan:   make(chan error, 1),
		unsubDone: make(chan struct{}),
	}

	go func(ctx context.Context) {
		defer close(sub.unsubDone)

		ticker := c.newTickerFn(ctx, c.Client)

		var fromBlockNum uint64
		var toBlockNum uint64

		for {
			select {
			case toBlockNum = <-ticker.Wait():
				if fromBlockNum == 0 {
					fromBlockNum = toBlockNum
				}

			case err, ok := <-ticker.Err():
				if !ok {
					errChan <- context.Canceled
					return
				}

				errChan <- err
				return

			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}

			if toBlockNum < fromBlockNum {
				continue
			}
			if toBlockNum > maxFilterLogWindow && toBlockNum-maxFilterLogWindow > fromBlockNum {
				toBlockNum = fromBlockNum + maxFilterLogWindow
			}

			q := filterQuery
			q.FromBlock = new(big.Int).SetUint64(fromBlockNum)
			q.ToBlock = new(big.Int).SetUint64(toBlockNum)

			logs, err := c.Client.FilterLogs(ctx, q)
			if err != nil {
				errChan <- err
				return
			}

			for _, log := range logs {
				select {
				case logChan <- log:
				case <-ctx.Done():
					errChan <- err
					return
				}
			}

			fromBlockNum = toBlockNum + 1
		}
	}(subCtx)

	return sub, nil
}

type clientSubscription struct {
	subCancel func()
	errChan   chan error
	unsubDone chan struct{}
}

func (s *clientSubscription) Unsubscribe() {
	s.subCancel()
	<-s.unsubDone

	close(s.errChan)
}

func (s *clientSubscription) Err() <-chan error {
	return s.errChan
}
