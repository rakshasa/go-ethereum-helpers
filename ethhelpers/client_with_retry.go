package ethhelpers

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// ClientWithRetry is a wrapper around the ethclient.Client that retries failed requests.
type ClientWithRetry interface {
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
}

type clientWithRetry struct {
	client       ClientWithRetry
	errorHandler func(context.Context, func(context.Context) error) error
}

func NewClientWithRetry(client ClientWithRetry, errorHandler func(context.Context, func(context.Context) error) error) ClientWithRetry {
	return &clientWithRetry{
		client:       client,
		errorHandler: errorHandler,
	}
}

func (c *clientWithRetry) BlockNumber(ctx context.Context) (r uint64, e error) {
	e = c.errorHandler(ctx, func(context.Context) (err error) {
		r, err = c.client.BlockNumber(ctx)
		return
	})
	return
}

func (c *clientWithRetry) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (r []types.Log, e error) {
	e = c.errorHandler(ctx, func(context.Context) (err error) {
		r, err = c.client.FilterLogs(ctx, q)
		return
	})
	return
}

// add method to client for SubscribeFilterlogs
func (c *clientWithRetry) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (sub ethereum.Subscription, e error) {
	e = c.errorHandler(ctx, func(context.Context) (err error) {
		sub, err = c.client.SubscribeFilterLogs(ctx, q, ch)
		return
	})
	return
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
			case errors.Is(err, syscall.ECONNRESET):
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
