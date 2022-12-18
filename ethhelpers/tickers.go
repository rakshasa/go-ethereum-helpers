package ethhelpers

import (
	"context"
	"fmt"
	"time"
)

// BlockNumberTicker is a ticker that emits block numbers.
type BlockNumberTicker interface {
	// Wait returns a channel that emits block numbers.
	Wait() <-chan uint64

	// Err returns a channel that emits errors that occur while waiting for block numbers.
	Err() <-chan error

	// CloneFromBlock creates a new ticker that starts from the given block number.
	CloneFromBlock(fromBlock uint64) BlockNumberTicker

	// Stop stops the ticker.
	Stop()
}

// blockNumberTicker is a generic handler for block number tickers.
type blockNumberTicker struct {
	interrupt      chan<- struct{}
	request        chan<- struct{}
	result         <-chan uint64
	errors         <-chan error
	cloneFromBlock func(uint64) *blockNumberTicker
	stop           func()
}

func (t *blockNumberTicker) Wait() <-chan uint64 {
	select {
	case t.interrupt <- struct{}{}:
	default:
	}

	if len(t.request) == 0 {
		t.request <- struct{}{}
	}

	return t.result
}

func (t *blockNumberTicker) Err() <-chan error {
	return t.errors
}

func (t *blockNumberTicker) CloneFromBlock(fromBlock uint64) BlockNumberTicker {
	return t.cloneFromBlock(fromBlock)
}

func (t *blockNumberTicker) Stop() {
	t.stop()
}

// TODO: Add discard duration, default to half of interval.
// TODO: Add max block interval and historic iteration options, these should wrap the PBNT Wait channel.
// TODO: Make fromBlock explicit, either it is set or we use the first value returned by BlockNumber.
// TODO: Note and verify that we start the ticker from the next block number.

// NewPeriodicBlockNumberTicker creates a new block number ticker that
// ticks at a fixed time interval, starting from the current block number.
func NewPeriodicBlockNumberTicker(ctx context.Context, client BlockNumberReader, interval time.Duration) BlockNumberTicker {
	ctx, stop := context.WithCancel(ctx)
	return newPeriodicBlockNumberTicker(ctx, stop, client, interval, 0)
}

// NewPeriodicBlockNumberTickerFromBlock creates a new block number ticker that
// ticks at a fixed time interval, starting from the given block number.
//
// The ticker continues to make BlockNumber request calls after calling Wait if
// fromBlock was not reached. Therefor the ticker should be manually stopped
// and/or not used with fromBlock values that are not imminient.
func NewPeriodicBlockNumberTickerFromBlock(ctx context.Context, client BlockNumberReader, interval time.Duration, fromBlock uint64) BlockNumberTicker {
	ctx, stop := context.WithCancel(ctx)
	return newPeriodicBlockNumberTicker(ctx, stop, client, interval, fromBlock)
}

type periodicBlockNumberTickerSource struct {
	client     BlockNumberReader
	interrupt  <-chan struct{}
	request    <-chan struct{}
	result     chan<- uint64
	errors     chan<- error
	windowSize uint64
}

func newPeriodicBlockNumberTicker(ctx context.Context, stop func(), client BlockNumberReader, interval time.Duration, fromBlock uint64) *blockNumberTicker {
	interrupt := make(chan struct{})
	request := make(chan struct{}, 1)
	result := make(chan uint64)
	errors := make(chan error, 1)

	t := &periodicBlockNumberTickerSource{
		client:    client,
		interrupt: interrupt,
		request:   request,
		result:    result,
		errors:    errors,
	}

	go t.start(ctx, interval, fromBlock)

	return &blockNumberTicker{
		interrupt: interrupt,
		request:   request,
		result:    result,
		errors:    errors,
		cloneFromBlock: func(fb uint64) *blockNumberTicker {
			return newPeriodicBlockNumberTicker(ctx, stop, client, interval, fb)
		},
		stop: stop,
	}
}

func (t *periodicBlockNumberTickerSource) start(ctx context.Context, interval time.Duration, fromBlock uint64) {
	defer close(t.errors)

	select {
	case <-t.request:
	case <-ctx.Done():
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		currentBlock, err := t.client.BlockNumber(ctx)
		if err != nil {
			t.errors <- err
			return
		}

		// Ensure the user doesn't overflow the uint64 block number if incremented twice.
		if currentBlock+2 < currentBlock {
			t.errors <- fmt.Errorf("block number overflow")
			return
		}

		if err := t.handle(ctx, &fromBlock, currentBlock, ticker.C); err != nil {
			t.errors <- err
			return
		}
	}
}

func (t *periodicBlockNumberTickerSource) handle(ctx context.Context, fromBlock *uint64, currentBlock uint64, tickerC <-chan time.Time) error {
	var requestC <-chan struct{}
	var resultC chan<- uint64

	if currentBlock >= *fromBlock {
		resultC = t.result

		// Make sure the request channel is empty before we attempt to send the result.
		select {
		case <-t.request:
		default:
		}
	}

	for {
		if requestC != nil && tickerC != nil {
			return fmt.Errorf("both the request and ticker channels are active")
		}

		select {
		case <-t.interrupt:
			// Avoid a race-condition where we are in a new tick and have a new
			// incoming request, but there's still unsent result available.
			if resultC != nil && tickerC == nil {
				// TODO: Only request again if +2 ticks or duration.
				resultC = nil
			}
			continue

		case resultC <- currentBlock:
			*fromBlock = currentBlock + 1

			// TODO: Do we keep consistent ticker time, or do we reset? Should be an option.
			resultC = nil
			continue

		case <-requestC:
			return nil

		case <-tickerC:
			if resultC != nil {
				// Last block number wasn't retrieved, so ignore the periodic ticker
				// and wait for a request before resuming.
				tickerC = nil
				requestC = t.request
				continue
			}

			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
