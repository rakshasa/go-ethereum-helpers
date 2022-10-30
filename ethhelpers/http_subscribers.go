package ethhelpers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type httpSubscription struct {
	subCancel func()
	errChan   chan error
	unsubDone chan struct{}
}

func (s *httpSubscription) Unsubscribe() {
	s.subCancel()
	<-s.unsubDone

	close(s.errChan)
}

func (s *httpSubscription) Err() <-chan error {
	return s.errChan
}

// The context argument cancels the RPC request that sets up the subscription
// but has no effect on the subscription after Subscribe has returned.
//
// Subscribers should be using the same underlying go-ethereum rpc client as the
// block number ticker to ensure there are no race-conditions.
//
// To conform to the SubscriberFilterLogs api the ticker should only return
// current, and not historic, block numbers.
func SubscribeFilterLogsWithHTTP(ctx context.Context, client FilterLogsReader, createTicker func(context.Context, uint64) BlockNumberTicker, filterQuery ethereum.FilterQuery, logChan chan<- types.Log) (ethereum.Subscription, error) {
	if logChan == nil {
		panic("channel given to SubscribeFilterLogs must not be nil")
	}

	subCtx, subCancel := context.WithCancel(context.Background())

	s := &httpSubscription{
		subCancel: subCancel,
		errChan:   make(chan error, 1),
		unsubDone: make(chan struct{}),
	}

	queryFromBlock, ok := BigIntAsUint64OrZeroIfNil(filterQuery.FromBlock)
	if !ok {
		return nil, fmt.Errorf("invalid FromBlock value")
	}

	ticker := createTicker(ctx, queryFromBlock)

	go func(ctx context.Context) {
		defer close(s.unsubDone)

		waitFn := func() (uint64, bool) {
			select {
			case num := <-ticker.Wait():
				return num, true

			case err := <-ticker.Err():
				if err == nil {
					s.errChan <- fmt.Errorf("block number ticker returned nil error")
					return 0, false
				}

				s.errChan <- err
				return 0, false

			case <-ctx.Done():
				s.errChan <- ctx.Err()
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

			q := filterQuery
			q.FromBlock = new(big.Int).SetUint64(fromBlock)
			q.ToBlock = new(big.Int).SetUint64(currentBlock)

			logs, err := client.FilterLogs(ctx, q)
			if err != nil {
				s.errChan <- err
				return
			}

			for _, log := range logs {
				select {
				case logChan <- log:
				case <-ctx.Done():
					s.errChan <- err
					return
				}
			}

			fromBlock = currentBlock + 1
		}
	}(subCtx)

	return s, nil
}
