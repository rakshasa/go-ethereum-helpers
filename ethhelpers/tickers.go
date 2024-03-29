package ethhelpers

import (
	"context"
	"fmt"
	"time"
)

// TODO: Add discard duration, default to half of interval.
// TODO: Add max distance between start block and fromBlock. Also sanity check current block returned.
// TODO: Replace with periodic ticker options config.

// TODO: Add a ResetToBlock method.
// TODO: Add a way to get the working block number.

// BlockNumberTicker is a ticker that emits block numbers.
type BlockNumberTicker interface {
	// Wait returns a channel that emits BlockNumber, and must be called before
	// each read.
	//
	// Reading from previously returned channels is not supported, and may in
	// some edge cases cause the new channel to skip a tick.
	//
	// To make sure the block number returned is up-to-date, call Wait rather
	// than reuse an unread channel.
	//
	// Truncated in the result is true if it was truncated by the window size.
	Wait() <-chan BlockNumber

	// Err returns a channel that emits errors that occur while waiting for block numbers.
	Err() <-chan error

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
	request chan blockNumberTickerRequest
	errors  <-chan error
	stop    func()
}

type blockNumberTickerRequest struct {
	result chan<- BlockNumber
}

func (t *blockNumberTicker) Wait() <-chan BlockNumber {
	ch := make(chan BlockNumber)

	for {
		select {
		case t.request <- blockNumberTickerRequest{ch}:
			return ch
		default:
		}

		// Discard previous request if it hasn't been read yet.
		select {
		case <-t.request:
		default:
		}

		// This should work:

		// select {
		// case t.request <- blockNumberTickerRequest{ch}:
		//      return ch
		// case <-t.request:
		//      // Discard previous request if it hasn't been read yet.
		// }
	}
}

func (t *blockNumberTicker) Err() <-chan error {
	return t.errors
}

func (t *blockNumberTicker) Stop() {
	t.stop()
}

type PeriodicBlockNumberTickerOptions struct {
	Client    BlockNumberReader
	FromBlock *uint64

	// Interval is the minimum time between API requests.
	//
	// Must be greater than 1ms.
	Interval time.Duration

	// The ticker continues to make BlockNumber request calls after calling Wait if
	// fromBlock was not reached. Therefor the ticker should be manually stopped
	// and/or not used with fromBlock values that are not imminient.
	WindowSize uint64
}

type periodicBlockNumberTickerSource struct {
	client  BlockNumberReader
	request <-chan blockNumberTickerRequest
	result  chan<- BlockNumber
	errors  chan<- error
	ticker  *time.Ticker

	fromBlock  uint64
	windowSize uint64
}

// NewPeriodicBlockNumberTicker creates a new block number ticker that
// ticks at a fixed time interval, starting from the current block number.
func NewPeriodicBlockNumberTicker(ctx context.Context, opts PeriodicBlockNumberTickerOptions) (BlockNumberTicker, error) {
	if opts.Client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if opts.Interval <= 1*time.Millisecond {
		return nil, fmt.Errorf("interval must be greater than 1ms")
	}

	ctx, stop := context.WithCancel(ctx)

	request := make(chan blockNumberTickerRequest, 1)
	errors := make(chan error, 1)

	t := &periodicBlockNumberTickerSource{
		client:     opts.Client,
		request:    request,
		errors:     errors,
		windowSize: opts.WindowSize,
	}

	go t.start(ctx, opts.Interval, opts.FromBlock)

	return &blockNumberTicker{
		request: request,
		errors:  errors,
		stop:    stop,
	}, nil
}

func (t *periodicBlockNumberTickerSource) start(ctx context.Context, interval time.Duration, initialFromBlock *uint64) {
	defer close(t.errors)

	select {
	case request := <-t.request:
		t.result = request.result
	case <-ctx.Done():
		t.errors <- ctx.Err()
		return
	}

	t.ticker = time.NewTicker(interval)
	defer t.ticker.Stop()

	currentBlock, err := t.client.BlockNumber(ctx)
	if err != nil {
		t.errors <- err
		return
	}

	if initialFromBlock == nil {
		t.fromBlock = currentBlock
	} else {
		t.fromBlock = *initialFromBlock
	}

	// TODO: Properly verify these sanity checks.

	if t.fromBlock+2 < t.fromBlock || t.fromBlock+t.windowSize < t.fromBlock {
		t.errors <- fmt.Errorf("from block number overflow")
		return
	}

	if t.windowSize+2 < t.windowSize {
		t.errors <- fmt.Errorf("window size overflow")
		return
	}

	for {
		if err := func() error {
			// Ensure the user doesn't overflow the uint64 block number if incremented twice.
			if currentBlock+2 < currentBlock || currentBlock+t.windowSize < currentBlock {
				return fmt.Errorf("block number overflow")
			}

			// TODO: Add option for an adaptive tick interval to avoid tickers silently
			// eating up api requests. Perhaps a Flush() method that reads requestC.
			if currentBlock < t.fromBlock {
				select {
				case <-t.ticker.C:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			// TODO: Add option to require time interval even for truncated results.
			// TODO: Re-request latest block number if ticker has passed.

			timestamp := time.Now()

			if t.windowSize != 0 && currentBlock > t.fromBlock+(t.windowSize-1) {
				if err := t.handleTruncated(ctx, currentBlock, timestamp); err != nil {
					return err
				}
			}

			if err = t.handleLatest(ctx, currentBlock, timestamp); err != nil {
				return err
			}

			return nil

		}(); err != nil {
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

// TODO: Add option to use interval for truncated results.

func (t *periodicBlockNumberTickerSource) handleTruncated(ctx context.Context, currentBlock uint64, timestamp time.Time) error {
	resultC := t.result

	for currentBlock > t.fromBlock+(t.windowSize-1) {
		select {
		case resultC <- BlockNumber{t.fromBlock + (t.windowSize - 1), timestamp, true}:
			t.fromBlock = t.fromBlock + (t.windowSize - 1) + 1
			t.result = nil
			resultC = nil

		case request := <-t.request:
			t.result = request.result
			resultC = request.result

		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// TODO: Not entirely correct...
	if resultC == nil {
		select {
		case request := <-t.request:
			t.result = request.result
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// TODO: If a new request is received and ticker has passed, do a new
	// request. This should be done when the for loop has written more than one
	// result.

	return nil
}

func (t *periodicBlockNumberTickerSource) handleLatest(ctx context.Context, currentBlock uint64, timestamp time.Time) error {
	// TODO: Reconsider this.
	if t.result == nil {
		panic("handleLatest: result channel is nil")
	}

	resultC := t.result
	tickerC := t.ticker.C

	shouldUpdate := false

	for {
		select {
		case request := <-t.request:
			t.result = request.result

			if tickerC == nil {
				return nil
			}

			// Check timestamp to see if we're close enough to the next tick.
			shouldUpdate = true

			if resultC != nil {
				resultC = request.result
			}

		case resultC <- BlockNumber{currentBlock, timestamp, false}:
			// TODO: Add option to reset ticker on result being received.
			t.fromBlock = currentBlock + 1
			t.result = nil
			resultC = nil

		case <-tickerC:
			tickerC = nil

			// TODO: Need to safeguard against slow api calls causing the ticker
			// to fire before the result is sent through the channel.
			if shouldUpdate {
				return nil
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
