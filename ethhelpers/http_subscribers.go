package ethhelpers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type httpSubscription struct {
	cancel func()
	err    chan error
	done   chan struct{}
}

func (s *httpSubscription) Unsubscribe() {
	s.cancel()
	<-s.done

	close(s.err)
}

func (s *httpSubscription) Err() <-chan error {
	return s.err
}

type HTTPSubscriberClient interface {
	BlockNumberReader
	FilterLogsReader
}

type HTTPSubscriberOptions struct {
	Client HTTPSubscriberClient

	// CreateContext returns a context that is used for the subscription, or
	// context.Background() if nil.
	CreateContext func() (context.Context, context.CancelFunc)

	// CreateTicker is a function that creates a block number ticker.
	//
	// The context is canceled when a subscription is unsubscribed or encounters
	// an error, or if the context passed to the function creating the
	// HTTPSubscriber function is canceled before returning a valid
	// subscription.
	//
	// The method is called only once per subscription.
	CreateTicker func(ctx context.Context, fromBlock uint64) (BlockNumberTicker, error)

	FilterQuery ethereum.FilterQuery
	Logs        chan<- types.Log
}

// The context argument cancels the RPC request that sets up the subscription
// but has no effect on the subscription after Subscribe has returned.
//
// Subscribers should be using the same underlying go-ethereum rpc client
// connection as the block number ticker to ensure there are no race-conditions.
//
// To conform to the go-ethereum SubscriberFilterLogs api the ticker should only
// return current, and not historic, block numbers.
//
// The current block number is requested before the subscription is returned.
func SubscribeFilterLogsWithHTTP(callerCtx context.Context, opts *HTTPSubscriberOptions) (ethereum.Subscription, error) {
	if opts.Client == nil {
		return nil, fmt.Errorf("opts.Client must be set")
	}
	if opts.CreateTicker == nil {
		return nil, fmt.Errorf("opts.CreateTicker must be set")
	}
	if opts.Logs == nil {
		return nil, fmt.Errorf("opts.Logs must be set")
	}

	subscriberCtx, cancel := func() (context.Context, context.CancelFunc) {
		if opts.CreateContext == nil {
			return context.WithCancel(context.Background())
		}

		return opts.CreateContext()
	}()

	ticker, err := func(ctx context.Context, done <-chan struct{}) (BlockNumberTicker, error) {
		ch := make(chan struct{})
		defer close(ch)

		go func() {
			select {
			case <-ch:
			case <-done:
				cancel()
			}
		}()

		currentBlock, err := opts.Client.BlockNumber(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get current block number: %w", err)
		}

		if opts.FilterQuery.FromBlock == nil {
			return opts.CreateTicker(ctx, currentBlock)
		}
		if !opts.FilterQuery.FromBlock.IsUint64() {
			return nil, fmt.Errorf("opts.FilterQuery.FromBlock is too large")
		}

		queryFromBlock := opts.FilterQuery.FromBlock.Uint64()

		if currentBlock < queryFromBlock {
			return opts.CreateTicker(ctx, queryFromBlock)
		} else {
			return opts.CreateTicker(ctx, currentBlock)
		}

	}(subscriberCtx, callerCtx.Done())
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create block ticker: %w", err)
	}

	s := &httpSubscription{
		cancel: cancel,
		err:    make(chan error, 1),
		done:   make(chan struct{}),
	}

	go func(ctx context.Context) {
		defer close(s.done)

		waitFn := func() (uint64, bool) {
			select {
			case bn, ok := <-ticker.Wait():
				if !ok {
					s.err <- fmt.Errorf("block ticker wait channel closed")
					return 0, false
				}

				return bn.BlockNumber, true

			case err, ok := <-ticker.Err():
				if !ok {
					s.err <- fmt.Errorf("block number ticker closed the error channel")
					return 0, false
				}
				if err == nil {
					s.err <- fmt.Errorf("block number ticker returned a nil error")
					return 0, false
				}

				s.err <- err
				return 0, false

			case <-ctx.Done():
				s.err <- ctx.Err()
				return 0, false
			}
		}

		fromBlock, ok := waitFn()
		if !ok {
			return
		}

		for {
			currentBlock, ok := waitFn()
			if !ok {
				return
			}
			if currentBlock < fromBlock {
				s.err <- fmt.Errorf("block number ticker returned a block number less than the from block")
				return
			}

			q := opts.FilterQuery
			q.FromBlock = new(big.Int).SetUint64(fromBlock)
			q.ToBlock = new(big.Int).SetUint64(currentBlock)

			logs, err := opts.Client.FilterLogs(ctx, q)
			if err != nil {
				s.err <- err
				return
			}

			for _, log := range logs {
				select {
				case opts.Logs <- log:
				case <-ctx.Done():
					s.err <- err
					return
				}
			}

			fromBlock = currentBlock + 1
		}
	}(subscriberCtx)

	return s, nil
}
