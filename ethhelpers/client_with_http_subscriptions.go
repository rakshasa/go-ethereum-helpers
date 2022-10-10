package ethhelpers

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type clientWithHTTPSubscriptions struct {
	Client
}

// TODO: Add options.

func NewClientWithHTTPSubscriptions(c Client) Client {
	return &clientWithHTTPSubscriptions{
		Client: c,
	}
}

// The context argument cancels the RPC request that sets up the subscription
// but has no effect on the subscription after Subscribe has returned.
func (c *clientWithHTTPSubscriptions) SubscribeFilterLogs(ctx context.Context, filterQuery ethereum.FilterQuery, logChan chan<- types.Log) (ethereum.Subscription, error) {
	if logChan == nil {
		panic("channel given to SubscribeFilterLogs must not be nil")
	}

	// var fromBlock

	if filterQuery.FromBlock == nil {
		blockNum, err := c.BlockNumber(ctx)
		if err != nil {
			return nil, err
		}

		filterQuery.FromBlock = new(big.Int).SetUint64(blockNum)
	}

	// filterQuery.ToBlock = new(big.Int)

	subCtx, subCancel := context.WithCancel(context.Background())
	errChan := make(chan error, 1)

	sub := &clientSubscription{
		subCancel: subCancel,
		errChan:   make(chan error, 1),
		unsubDone: make(chan struct{}),
	}

	go func(ctx context.Context) {
		defer close(sub.unsubDone)

		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
			}

			// TODO: Add a head block number ticker type helper for use with all
			// block polling queries.
			blockNum, err := c.BlockNumber(ctx)
			if err != nil {
				errChan <- err
				return
			}

			// TOOD: Change to always use uint64 values.
			q := filterQuery
			q.ToBlock = new(big.Int).SetUint64(blockNum)

			if q.ToBlock.Cmp(q.FromBlock) < 0 {
				time.Sleep(3 * time.Second)
				continue
			}

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

			filterQuery.FromBlock.Add(q.ToBlock, big.NewInt(1))

			// TODO: Replace with ticker.
			time.Sleep(3 * time.Second)
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
