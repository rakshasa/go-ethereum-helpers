package ethhelpers

import (
	"context"
	"fmt"
	"time"
)

type BlockNumberTicker interface {
	Wait() <-chan uint64
	Err() <-chan error
	Reset(fromBlock uint64)
	Stop()
}

type blockNumberTicker struct {
	interruptChan chan<- struct{}
	requestChan   chan<- struct{}
	resultChan    <-chan uint64
	errChan       <-chan error
	reset         func(uint64) blockNumberTicker
	stop          func()
}

func (t *blockNumberTicker) Wait() <-chan uint64 {
	select {
	case t.interruptChan <- struct{}{}:
	default:
	}

	if len(t.requestChan) == 0 {
		t.requestChan <- struct{}{}
	}

	return t.resultChan
}

func (t *blockNumberTicker) Err() <-chan error {
	return t.errChan
}

func (t *blockNumberTicker) Reset(fromBlock uint64) {
	*t = t.reset(fromBlock)
}

func (t *blockNumberTicker) Stop() {
	t.stop()
}

// Add discard duration.

// TODO: Add max block interval and historic iteration options, these should wrap the PBNT Wait channel.

type periodicBlockNumberTicker struct {
	client        BlockNumberReader
	interruptChan <-chan struct{}
	requestChan   <-chan struct{}
	resultChan    chan<- uint64
	errChan       chan<- error
}

func NewPeriodicBlockNumberTicker(ctx context.Context, client BlockNumberReader, fromBlock uint64, interval time.Duration) BlockNumberTicker {
	ctx, stop := context.WithCancel(ctx)

	t := resetPeriodicBlockNumberTicker(ctx, stop, client, fromBlock, interval)
	return &t
}

func FactoryForPeriodicBlockNumberTicker(client BlockNumberReader, fromBlock uint64, interval time.Duration) func(context.Context) BlockNumberTicker {
	return func(ctx context.Context) BlockNumberTicker {
		return NewPeriodicBlockNumberTicker(ctx, client, fromBlock, interval)
	}
}

func FactoryForPeriodicBlockNumberTickerWithFromBlock(client BlockNumberReader, interval time.Duration) func(context.Context, uint64) BlockNumberTicker {
	return func(ctx context.Context, fromBlock uint64) BlockNumberTicker {
		return NewPeriodicBlockNumberTicker(ctx, client, fromBlock, interval)
	}
}

func resetPeriodicBlockNumberTicker(ctx context.Context, stop func(), client BlockNumberReader, fromBlock uint64, interval time.Duration) blockNumberTicker {
	interruptChan := make(chan struct{})
	requestChan := make(chan struct{}, 1)
	resultChan := make(chan uint64)
	errChan := make(chan error, 1)

	t := &periodicBlockNumberTicker{
		client:        client,
		interruptChan: interruptChan,
		requestChan:   requestChan,
		resultChan:    resultChan,
		errChan:       errChan,
	}

	go t.start(ctx, fromBlock, interval)

	return blockNumberTicker{
		interruptChan: interruptChan,
		requestChan:   requestChan,
		resultChan:    resultChan,
		errChan:       errChan,
		reset: func(newFromBlock uint64) blockNumberTicker {
			return resetPeriodicBlockNumberTicker(ctx, stop, client, newFromBlock, interval)
		},
		stop: stop,
	}
}

func (t *periodicBlockNumberTicker) start(ctx context.Context, fromBlock uint64, interval time.Duration) {
	defer close(t.errChan)

	select {
	case <-t.requestChan:
	case <-ctx.Done():
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		currentBlock, err := t.client.BlockNumber(ctx)
		if err != nil {
			// TODO: Allow passing custom error handlers.
			// TODO: Add a way to reset
			t.errChan <- err
			return
		}

		// Ensure the user doesn't overflow the uint64 block number if incremented twice.
		if currentBlock+2 < currentBlock {
			t.errChan <- fmt.Errorf("block number overflow")
			return
		}

		if err := t.handle(ctx, &fromBlock, currentBlock, ticker.C); err != nil {
			t.errChan <- err
			return
		}
	}
}

func (t *periodicBlockNumberTicker) handle(ctx context.Context, fromBlock *uint64, currentBlock uint64, tickerC <-chan time.Time) error {
	var requestC <-chan struct{}
	var resultC chan<- uint64

	if currentBlock >= *fromBlock {
		resultC = t.resultChan

		select {
		case <-t.requestChan:
		default:
		}
	}

	for {
		if requestC != nil && tickerC != nil {
			return fmt.Errorf("both the request and ticker channels are active")
		}

		select {
		case <-t.interruptChan:
			// Avoid a race-condition where we are in a new tick and have a new
			// incoming request, but there's still an unsent result available.
			if resultC != nil && tickerC == nil {
				// TODO: Only request again if +2 ticks or duration.
				resultC = nil
			}
			continue

		case <-requestC:
			// check if we should discard due to discard timeout, otherwise send old value
			return nil

		case resultC <- currentBlock:
			*fromBlock = currentBlock + 1

			// TODO: Do we keep consistent ticker time, or do we reset? Should be an option.
			resultC = nil
			continue

		case <-tickerC:
			if resultC == nil {
				return nil
			}

			tickerC = nil
			requestC = t.requestChan
			continue

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
