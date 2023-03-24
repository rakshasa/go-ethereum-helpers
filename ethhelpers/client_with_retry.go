package ethhelpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func NewClientWithRetry(client Client, errorHandler func(context.Context, func(context.Context) error) error) Client {
	return NewClientWithDefaultHandler(func(ctx context.Context, caller ClientCaller) error {
		return errorHandler(ctx, func(ctx context.Context) error {
			return caller.Call(ctx, client)
		})
	})
}

// TODO: Add retry options struct.
// TODO: Have different handling of calls that are expected to be valid, vs. might be invalid.

func RetryIfTemporaryError(unknownError func(context.Context, error) error) func(context.Context, func(context.Context) error) error {
	return func(ctx context.Context, fn func(context.Context) error) error {
		for {
			err := fn(ctx)
			if err == nil {
				return nil
			}

			var newErr error
			var rpcErr rpc.Error

			switch {
			case errors.Is(err, context.Canceled):
				return err
			case errors.Is(err, context.DeadlineExceeded):
				return err
			case errors.Is(err, syscall.ECONNRESET) || os.IsTimeout(err):
				// TODO: Use an temporary error handler.
				time.Sleep(1 * time.Second)
				continue
			case errors.As(err, &rpcErr):
				// TODO: Retry depending on the error.
				newErr = unknownError(ctx, fmt.Errorf("ethereum client call failed, rpc error: %w", err))
			default:
				newErr = unknownError(ctx, fmt.Errorf("ethereum client call failed, unknown error: %w", err))
			}

			if newErr != nil {
				return newErr
			}

			// TODO: Use back-off function.
			time.Sleep(3 * time.Second)
		}
	}
}
