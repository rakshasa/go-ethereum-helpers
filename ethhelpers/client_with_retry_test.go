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

func TestClientWithRetry_RetryIfTemporaryError(t *testing.T) {
	type testArgs struct {
		ctx    context.Context
		client ethhelpers.ClientWithRetry
		mock   *mock.Mock
	}

	tests := []struct {
		name   string
		client func(*testing.T, ethhelpers.ClientWithRetry) ethhelpers.ClientWithRetry
		fn     func(*testing.T, testArgs)
	}{
		{
			name: "BlockNumber immediately returns",
			fn: func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", args.ctx).Return(uint64(123), nil).Once()

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.NoError(t, err)
				assert.Equal(t, uint64(123), blockNumber)
			},
		}, {
			name: "BlockNumber returns after an initial connection timeout",
			fn: func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", args.ctx).Return(uint64(0), syscall.ECONNRESET).Once().After(10 * time.Millisecond)
				args.mock.On("BlockNumber", args.ctx).Return(uint64(123), nil).Once().After(10 * time.Millisecond)

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.NoError(t, err)
				assert.Equal(t, uint64(123), blockNumber)
			},
		}, {
			name: "BlockNumber returns after context is canceled",
			fn: func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", args.ctx).Return(uint64(0), context.Canceled).Once().After(10 * time.Millisecond)

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.ErrorIs(t, err, context.Canceled)
				assert.Equal(t, uint64(0), blockNumber)
			},
		}, {
			name: "BlockNumber returns after context is deadline exceeded",
			fn: func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", args.ctx).Return(uint64(0), context.DeadlineExceeded).Once().After(10 * time.Millisecond)

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.ErrorIs(t, err, context.DeadlineExceeded)
				assert.Equal(t, uint64(0), blockNumber)
			},
		}, {
			name: "FilterLogs immediately returns",
			fn: func(t *testing.T, args testArgs) {
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
			t.Parallel()

			ctx := context.Background()
			client := ethtesting.NewClientWithMock()

			test.fn(t, testArgs{
				ctx,
				ethhelpers.NewClientWithRetry(client, ethhelpers.RetryIfTemporaryError(unknownErrorWithAssertFail(t))),
				client.Mock(),
			})

			client.Mock().AssertExpectations(t)
		})
	}
}
