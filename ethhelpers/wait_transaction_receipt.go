package ethhelpers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type WaitForTransactionReceiptOptions struct {
	// TODO: Add TransactionReaderFromContext, use it if Client is nil.
	Client ethereum.TransactionReader
	TxHash common.Hash

	// ErrorHandler is called when Client.TransactionReceipt returns a non-nil
	// error. The handler is not called if the main context was canceled.
	//
	// If the handler returns a non-nil error then the error is sent to the
	// result channel and no further attempts are made.
	ErrorHandler func(txHash common.Hash, err error) error
}

// TODO: Move to types.
type ReceiptOrError struct {
	Receipt *types.Receipt
	Error   error
}

// WaitForTransactionReceipt
//
// The result channel always sends either the receipt or error, and does not close.
//
// The caller must ensure that either the context is canceled or the returned
// cancel function is called.
func WaitForTransactionReceipt(ctx context.Context, options WaitForTransactionReceiptOptions) (<-chan ReceiptOrError, func()) {
	// TODO: Add a FooWithCancel variant.
	ctx, cancel := context.WithCancel(ctx)

	resultChan := make(chan ReceiptOrError, 1)

	go func() {
		// TODO: Increase timeout per attempt, add options for timeouts.
		// TODO: Add defaults based on chain config from context..
		// TODO: Add an option to use "sub, err := client.SubscribeNewHead(context.Background(), headers)"
		// TODO: Return error if txHash is zero value.

		ticker := time.NewTicker(time.Duration(3) * time.Second)
		defer ticker.Stop()

		for {
			checkFn := func() (*types.Receipt, error) {
				ctx, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Minute)
				defer cancel()

				return options.Client.TransactionReceipt(ctx, options.TxHash)
			}

			receipt, err := checkFn()
			if err == nil {
				resultChan <- ReceiptOrError{receipt, nil}
				return
			}

			// Always cancel if the parent context was cancelled.
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				select {
				case <-ctx.Done():
					resultChan <- ReceiptOrError{nil, ctx.Err()}
					return
				default:
				}
			}

			if err := options.ErrorHandler(options.TxHash, err); err != nil {
				resultChan <- ReceiptOrError{nil, err}
				return
			}

			// TODO: Add custom ticker and add testing.

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				resultChan <- ReceiptOrError{nil, ctx.Err()}
				return
			}
		}
	}()

	return resultChan, cancel
}

type WaitTransactionReceipts struct {
	mu sync.Mutex

	ctx    context.Context
	cancel func()
	count  int

	// TODO: The cancel function is probably unnessesary.
	addFn       func(context.Context, common.Hash) (<-chan ReceiptOrError, func())
	collectChan chan ReceiptOrError
	resultChan  chan ReceiptOrError
}

func NewWaitTransactionReceipts(ctx context.Context, addFn func(context.Context, common.Hash) (<-chan ReceiptOrError, func())) *WaitTransactionReceipts {
	ctx, cancel := context.WithCancel(ctx)

	w := &WaitTransactionReceipts{
		ctx:    ctx,
		cancel: cancel,

		addFn:       addFn,
		collectChan: make(chan ReceiptOrError, 16),
		resultChan:  make(chan ReceiptOrError, 1),
	}

	pickErrorFn := func(currentError, newError error) error {
		switch {
		case currentError == nil:
			return newError
		case currentError == context.Canceled:
			return newError
		case newError == context.Canceled:
			return currentError
		case newError == context.DeadlineExceeded:
			return currentError
		default:
			return newError
		}
	}

	go func() {
		defer cancel()

		var currentError error

		for {
			select {
			case result := <-w.collectChan:
				w.mu.Lock()
				w.count--

				if result.Error == nil {
					w.resultChan <- result

					cancel()
					w.mu.Unlock()
					return
				}

				// TODO: Need to improve the collection of errors to make it
				// prefer the lastest txhash, perhaps add an index.
				currentError = pickErrorFn(currentError, result.Error)

				if w.count != 0 {
					if errors.Is(currentError, context.Canceled) || errors.Is(currentError, context.DeadlineExceeded) {
						w.mu.Unlock()
						continue
					}
				}

				w.resultChan <- ReceiptOrError{Error: currentError}

				cancel()
				w.mu.Unlock()
				return

			case <-ctx.Done():
				w.resultChan <- ReceiptOrError{Error: ctx.Err()}
				return
			}
		}
	}()

	return w
}

// TODO: Add different different WatchFor* functions that allow us to use either
// pooling TransactionReceipt or a websocket subscription.
func (w *WaitTransactionReceipts) Add(txHash common.Hash) {
	w.mu.Lock()
	defer w.mu.Unlock()

	ch, cancel := w.addFn(w.ctx, txHash)

	go func() {
		w.collectChan <- (<-ch)
		cancel()
	}()

	w.count++
}

// Only read channel once.
func (w *WaitTransactionReceipts) Result() <-chan ReceiptOrError {
	return w.resultChan
}

func (w *WaitTransactionReceipts) Stop() {
	w.cancel()
}
