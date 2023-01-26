package ethhelpers

import (
	"context"
	"fmt"
	"time"
)

// BlockNumberTicker is a ticker that emits block numbers.
type BlockNumberTicker interface {
	// Wait returns a channel that emits BlockNumber, and must be called at
	// least once per result read.
	//
	// Reading from old channels has undefined behavior, and no reads should be
	// ongoing when calling Wait.
	//
	// The returned channel is only guaranteed to return a valid result once,
	// reading from the channel multiple times or from multiple channels
	// returned by Wait has undefined behavior.
	//
	// Truncated in the result is true if it was truncated by the window size.
	Wait() <-chan BlockNumber

	// Err returns a channel that emits errors that occur while waiting for block numbers.
	Err() <-chan error

	// CloneFromBlock creates a new ticker that starts from the given block number.
	CloneFromBlock(fromBlock uint64) BlockNumberTicker

	// Stop stops the ticker.
	Stop()
}

type BlockNumber struct {
	BlockNumber uint64
	Timestamp   time.Time
	Truncated   bool
}

// blockNumberTicker is a generic handler for block number tickers.
type blockNumberTicker struct {
	request        chan<- struct{}
	result         <-chan BlockNumber
	errors         <-chan error
	cloneFromBlock func(uint64) *blockNumberTicker
	stop           func()
}

func (t *blockNumberTicker) Wait() <-chan BlockNumber {
	select {
	case t.request <- struct{}{}:
	default:
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

// NewPeriodicBlockNumberTicker creates a new block number ticker that
// ticks at a fixed time interval, starting from the current block number.
func NewPeriodicBlockNumberTicker(ctx context.Context, client BlockNumberReader, interval time.Duration) BlockNumberTicker {
	ctx, stop := context.WithCancel(ctx)
	return newPeriodicBlockNumberTicker(ctx, stop, client, interval, nil, 0)
}

func NewPeriodicBlockNumberTickerWithWindowSize(ctx context.Context, client BlockNumberReader, interval time.Duration, windowSize uint64) BlockNumberTicker {
	ctx, stop := context.WithCancel(ctx)
	return newPeriodicBlockNumberTicker(ctx, stop, client, interval, nil, windowSize)
}

// NewPeriodicBlockNumberTickerFromBlock creates a new block number ticker that
// ticks at a fixed time interval, starting from the given block number.
//
// The ticker continues to make BlockNumber request calls after calling Wait if
// fromBlock was not reached. Therefor the ticker should be manually stopped
// and/or not used with fromBlock values that are not imminient.
func NewPeriodicBlockNumberTickerFromBlock(ctx context.Context, client BlockNumberReader, interval time.Duration, fromBlock uint64) BlockNumberTicker {
	ctx, stop := context.WithCancel(ctx)
	return newPeriodicBlockNumberTicker(ctx, stop, client, interval, &fromBlock, 0)
}

type periodicBlockNumberTickerSource struct {
	client     BlockNumberReader
	request    <-chan struct{}
	result     chan<- BlockNumber
	errors     chan<- error
	windowSize uint64
}

func newPeriodicBlockNumberTicker(ctx context.Context, stop func(), client BlockNumberReader, interval time.Duration, fromBlock *uint64, windowSize uint64) *blockNumberTicker {
	request := make(chan struct{}, 1)
	result := make(chan BlockNumber)
	errors := make(chan error, 1)

	t := &periodicBlockNumberTickerSource{
		client:     client,
		request:    request,
		result:     result,
		errors:     errors,
		windowSize: windowSize,
	}

	go t.start(ctx, interval, fromBlock)

	return &blockNumberTicker{
		request: request,
		result:  result,
		errors:  errors,
		cloneFromBlock: func(fb uint64) *blockNumberTicker {
			return newPeriodicBlockNumberTicker(ctx, stop, client, interval, &fb, windowSize)
		},
		stop: stop,
	}
}

func (t *periodicBlockNumberTickerSource) start(ctx context.Context, interval time.Duration, initialFromBlock *uint64) {
	defer close(t.errors)

	select {
	case <-t.request:
	case <-ctx.Done():
		t.errors <- ctx.Err()
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	currentBlock, err := t.client.BlockNumber(ctx)
	if err != nil {
		t.errors <- err
		return
	}

	var fromBlock uint64

	if initialFromBlock == nil {
		fromBlock = currentBlock
	} else {
		fromBlock = *initialFromBlock
	}

	if fromBlock+2 < fromBlock || fromBlock+t.windowSize < fromBlock {
		t.errors <- fmt.Errorf("from block number overflow")
		return
	}

	for {
		// Ensure the user doesn't overflow the uint64 block number if incremented twice.
		if currentBlock+2 < currentBlock || currentBlock+t.windowSize < currentBlock {
			t.errors <- fmt.Errorf("block number overflow")
			return
		}

		// Make it loop until within windowSize, and then re-request the latest if ticker has passed.
		// Needs to be a separate handler.

		// TODO: Add option to require time interval even for truncated results.
		// TODO: Re-request latest block number if ticker has passed.

		// TODO: Consider rewrting this to have 'request' return a result
		// channel. Or rather, request passes a result channel to the ticker.

		timestamp := time.Now()

		if t.windowSize != 0 {
			if fromBlock, err = t.handleTruncated(ctx, fromBlock, currentBlock, timestamp); err != nil {
				t.errors <- err
				return
			}
		}

		if fromBlock, err = t.handleLatest(ctx, fromBlock, currentBlock, timestamp, ticker.C); err != nil {
			t.errors <- err
			return
		}

		// TODO: Mock client isn't canceling on context cancel.

		if currentBlock, err = t.client.BlockNumber(ctx); err != nil {
			t.errors <- err
			return
		}
	}
}

func (t *periodicBlockNumberTickerSource) handleTruncated(ctx context.Context, fromBlock, currentBlock uint64, timestamp time.Time) (uint64, error) {
	for fromBlock+(t.windowSize-1) < currentBlock {
		select {
		case t.result <- BlockNumber{fromBlock + (t.windowSize - 1), timestamp, true}:
			// TODO: Do we keep consistent ticker time, or do we reset? Should be an option.
			fromBlock = fromBlock + (t.windowSize - 1) + 1
		case <-ctx.Done():
			return fromBlock, ctx.Err()
		}
	}

	// TODO: If a new request is received and ticker has passed, do a new request.

	return fromBlock, nil
}

func (t *periodicBlockNumberTickerSource) handleLatest(ctx context.Context, fromBlock, currentBlock uint64, timestamp time.Time, tickerC <-chan time.Time) (uint64, error) {
	// We need to catch up to fromBlock, so just wait for the next tick.
	//
	// TODO: Add option for an adaptive tick interval to avoid tickers silently
	// eating up api requests. Perhaps a Flush() method that reads requestC.
	if currentBlock < fromBlock {
		select {
		case <-tickerC:
			return fromBlock, nil
		case <-ctx.Done():
			return fromBlock, ctx.Err()
		}
	}

	select {
	case <-t.request:
	default:
	}

	resultC := t.result

	for {
		select {
		case resultC <- BlockNumber{currentBlock, timestamp, false}:
			// TODO: Do we keep consistent ticker time, or do we reset? Should be an option.
			fromBlock = currentBlock + 1
			resultC = nil

		case <-tickerC:
			if resultC != nil {
				select {
				case <-t.request:
				default:
				}
			}

			for {
				select {
				case resultC <- BlockNumber{currentBlock, timestamp, false}:
					fromBlock = currentBlock + 1
					resultC = nil

				case <-t.request:
					// TODO: If we have race condition here we can end up passing a stale result.
					return fromBlock, nil

				case <-ctx.Done():
					return fromBlock, ctx.Err()
				}
			}

		case <-ctx.Done():
			return fromBlock, ctx.Err()
		}
	}
}
