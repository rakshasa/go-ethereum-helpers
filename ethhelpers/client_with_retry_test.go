package ethhelpers_test

import (
	"context"
	"math/big"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func unknownErrorWithAssertFail(t *testing.T) func(context.Context, error) error {
	return func(ctx context.Context, err error) error {
		assert.Failf(t, "should not be called", "%v", err)
		return context.Canceled
	}
}

func TestClientWithRetry_RetryIfTemporaryError(t *testing.T) {
	type testArgs struct {
		ctx    context.Context
		client ethhelpers.ClientWithRetry
		mock   *mock.Mock
	}

	defaultClient := func(t *testing.T, c ethhelpers.ClientWithRetry) ethhelpers.ClientWithRetry {
		return ethhelpers.NewClientWithRetry(c, ethhelpers.RetryIfTemporaryError(unknownErrorWithAssertFail(t)))
	}

	tests := []struct {
		name   string
		client func(*testing.T, ethhelpers.ClientWithRetry) ethhelpers.ClientWithRetry
		mock   func(*testing.T, testArgs)
	}{
		{
			name:   "BlockNumber immediately returns",
			client: defaultClient,
			mock: func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", args.ctx).Return(uint64(123), nil).Once()

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.NoError(t, err)
				assert.Equal(t, uint64(123), blockNumber)
			},
		}, {
			name:   "BlockNumber returns after an initial connection timeout",
			client: defaultClient,
			mock: func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", args.ctx).Return(uint64(0), syscall.ECONNRESET).Once().After(10 * time.Millisecond)
				args.mock.On("BlockNumber", args.ctx).Return(uint64(123), nil).Once().After(10 * time.Millisecond)

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.NoError(t, err)
				assert.Equal(t, uint64(123), blockNumber)
			},
		}, {
			name:   "BlockNumber returns after context is canceled",
			client: defaultClient,
			mock: func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", args.ctx).Return(uint64(0), context.Canceled).Once().After(10 * time.Millisecond)

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.ErrorIs(t, err, context.Canceled)
				assert.Equal(t, uint64(0), blockNumber)
			},
		}, {
			name:   "BlockNumber returns after context is deadline exceeded",
			client: defaultClient,
			mock: func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", args.ctx).Return(uint64(0), context.DeadlineExceeded).Once().After(10 * time.Millisecond)

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, uint64(0), blockNumber)
			},
		}, {
			name:   "FilterLogs immediately returns",
			client: defaultClient,
			mock: func(t *testing.T, args testArgs) {
				q := ethereum.FilterQuery{
					FromBlock: big.NewInt(1),
				}
				fl := []types.Log{types.Log{}}

				args.mock.On("FilterLogs", args.ctx, q).Return(fl, nil).Once()

				logs, err := args.client.FilterLogs(args.ctx, q)
				assert.NoError(t, err)
				assert.Equal(t, fl, logs)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			client, mock := ethtesting.NewClientWithMock()
			client = test.client(t, client)

			test.mock(t, testArgs{ctx, client, mock})

			mock.AssertExpectations(t)
		})
	}
}
