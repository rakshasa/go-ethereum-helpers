package ethhelpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type WaitForTransactionReceiptOptions struct {
	Client ethereum.TransactionReader
	TxHash common.Hash

	// ErrorHandler is called when Client.TransactionReceipt returns a non-nil
	// error. The handler is not called if the main context was canceled.
	//
	// If the handler returns a non-nil error then the error is sent to the
	// result channel and no further attempts are made.
	ErrorHandler func(err error) error
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
func WaitForTransactionReceipt(ctx context.Context, options WaitForTransactionReceiptOptions) (<-chan ReceiptOrError, func()) {
	ctx, cancel := context.WithCancel(ctx)

	resultChan := make(chan ReceiptOrError, 1)

	go func() {
		// TODO: Increase timeout per attempt, add options for timeouts.
		// TODO: Add defaults based on chain config from context..
		// TODO: Add an option to use "sub, err := client.SubscribeNewHead(context.Background(), headers)"

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

			if err := options.ErrorHandler(err); err != nil {
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

// DefaultErrorHandlerWithLogger is used by e.g. WaitForTransactionReceipt to
// decide if the error is a temporary connection / receipt not found issue and
// it should retry.
//
// This is a work-in-progress.
func DefaultErrorHandlerWithMessages(msgHandler func(msg string)) func(err error) error {
	return func(err error) error {
		var rpcErr rpc.Error

		switch {
		case errors.As(err, &rpcErr):
			// TODO: We should have optional permissive error handling code for
			// temporary errors.
			if err.Error() == "request failed or timed out" {
				msgHandler(fmt.Sprintf("waiting for transaction receipt, temporary error : %d : %v", rpcErr.ErrorCode(), rpcErr))
				break
			}

			msgHandler(fmt.Sprintf("waiting for transaction receipt, unknown rpc error : %d : %v", rpcErr.ErrorCode(), rpcErr))
			return err

		case errors.Is(err, context.Canceled):
			msgHandler(fmt.Sprintf("waiting for transaction receipt, context canceled"))

		case errors.Is(err, context.DeadlineExceeded):
			msgHandler(fmt.Sprintf("waiting for transaction receipt, context deadline"))

		case errors.Is(err, ethereum.NotFound):
			msgHandler(fmt.Sprintf("waiting for transaction receipt, not found"))

		case errors.Is(err, syscall.ECONNRESET):
			msgHandler(fmt.Sprintf("waiting for transaction receipt, connection reset by peer"))

		case os.IsTimeout(err):
			msgHandler(fmt.Sprintf("waiting for transaction receipt, timeout"))

		// TODO: Check non-temporary errors.
		default:
			msgHandler(fmt.Sprintf("waiting for transaction receipt, unknown error : %v", err))
			return err
		}

		return nil
	}
}
