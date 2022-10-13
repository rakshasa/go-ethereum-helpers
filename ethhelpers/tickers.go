package ethhelpers

import (
	"context"
	"errors"
	"time"
)

type BlockNumberTicker interface {
	Wait() <-chan uint64
	Err() <-chan error
	Stop()
}

type blockNumberTickerWithDuration struct {
	waitChan <-chan uint64
	errChan  <-chan error
	stop     func()
}

// TODO: Figure out how to ensure the latest block number is sent.
//
// Have different behaviors, e.g. continous nums, lazy request.

func NewBlockNumberTickerWithDuration(ctx context.Context, client BlockNumberReader, d time.Duration) BlockNumberTicker {
	ctx, stop := context.WithCancel(ctx)

	waitChan := make(chan uint64, 1)
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		ticker := time.NewTicker(d)
		defer ticker.Stop()

		var lastBlockNum uint64

		for {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}

			blockNum, err := client.BlockNumber(ctx)
			if err != nil {
				// TODO: Canceling context should send error, not close error channel.
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return
				}

				// TODO: Allow passing custom error handlers.
				select {
				case errChan <- err:
				default:
				}

				continue
			}

			if blockNum == lastBlockNum {
				continue
			}

			select {
			case waitChan <- blockNum:
				lastBlockNum = blockNum
			case <-ctx.Done():
				return
			}
		}
	}()

	return &blockNumberTickerWithDuration{
		waitChan: waitChan,
		errChan:  errChan,
		stop:     stop,
	}
}

func FactoryForNewBlockNumberTickerWithDuration(d time.Duration) func(context.Context, BlockNumberReader) BlockNumberTicker {
	return func(ctx context.Context, client BlockNumberReader) BlockNumberTicker {
		return NewBlockNumberTickerWithDuration(ctx, client, d)
	}
}

func (t *blockNumberTickerWithDuration) Wait() <-chan uint64 {
	return t.waitChan
}

func (t *blockNumberTickerWithDuration) Err() <-chan error {
	return t.errChan
}

func (t *blockNumberTickerWithDuration) Stop() {
	t.stop()
}
