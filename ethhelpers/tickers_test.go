package ethhelpers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTickers_NewPeriodicBlockNumberTicker(t *testing.T) {
	type testArgs struct {
		ctx    context.Context
		client ethtesting.ClientWithMock
		mock   *mock.Mock
		ticker ethhelpers.BlockNumberTicker
	}

	// TODO: Add tester client that simulates errors.
	// TODO: Add test that wait with <1/4s sleep.
	// TODO: Test with repeated wait calls.
	// TODO: Test errors from client.BlockNumber(ctx)
	// TODO: Test sim starting with non-zero block number.

	tests := []struct {
		name string
		fn   func(t *testing.T, args testArgs)
	}{
		{
			"has empty channels immediately after wait call",
			func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(10), nil).Once().After(200 * time.Millisecond)

				ch := args.ticker.Wait()
				time.Sleep(50 * time.Millisecond)

				assert.True(t, emptyUint64Channel(ch))
				assert.True(t, emptyErrorChannel(args.ticker.Err()))
			},
		}, {
			"single request, read with timeout",
			func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(10), nil).Once().After(20 * time.Millisecond)

				ch := args.ticker.Wait()
				assert.True(t, emptyUint64Channel(ch))

				bn, ok := readUint64FromChanWithTimeout(ch, 50*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, uint64(10), bn)

				assert.True(t, emptyErrorChannel(args.ticker.Err()))
			},
		}, {
			"single request, read after sleeping",
			func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(10), nil).Once().After(20 * time.Millisecond)

				ch := args.ticker.Wait()
				assert.True(t, emptyUint64Channel(ch))

				time.Sleep(50 * time.Millisecond)

				bn, ok := readUint64FromChan(ch)
				assert.True(t, ok)
				assert.Equal(t, uint64(10), bn)

				assert.True(t, emptyErrorChannel(args.ticker.Err()))
			},
		}, {
			"single request, read after sleeping beyond tick",
			func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(10), nil).Once().After(20 * time.Millisecond)

				ch := args.ticker.Wait()

				time.Sleep(150 * time.Millisecond)

				bn, ok := readUint64FromChan(ch)
				assert.True(t, ok)
				assert.Equal(t, uint64(10), bn)

				assert.True(t, emptyErrorChannel(args.ticker.Err()))
			},
		}, {
			"double request, read next immediately, with new block number",
			func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(10), nil).Once().After(10 * time.Millisecond)
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(11), nil).Once().After(10 * time.Millisecond)

				ch := args.ticker.Wait()

				bn, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, uint64(10), bn)

				ch = args.ticker.Wait()

				bn, ok = readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.False(t, ok)
				assert.Equal(t, uint64(0), bn)

				time.Sleep(100 * time.Millisecond)

				bn, ok = readUint64FromChan(ch)
				assert.True(t, ok)
				assert.Equal(t, uint64(11), bn)

				assert.True(t, emptyErrorChannel(args.ticker.Err()))
			},
		}, {
			"double request, read after sleeping beyond tick, with new block number",
			func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(10), nil).Once().After(10 * time.Millisecond)
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(11), nil).Once().After(10 * time.Millisecond)

				ch := args.ticker.Wait()

				bn, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, uint64(10), bn)

				time.Sleep(100 * time.Millisecond)

				ch = args.ticker.Wait()

				bn, ok = readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, uint64(11), bn)

				assert.True(t, emptyErrorChannel(args.ticker.Err()))
			},
		}, {
			"double request, read next immediately, with same block number",
			func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(10), nil).Once().After(10 * time.Millisecond)
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(10), nil).Once().After(10 * time.Millisecond)

				ch := args.ticker.Wait()

				bn, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, uint64(10), bn)

				ch = args.ticker.Wait()

				_, ok = readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.False(t, ok)

				time.Sleep(100 * time.Millisecond)

				_, ok = readUint64FromChan(ch)
				assert.False(t, ok)

				assert.True(t, emptyErrorChannel(args.ticker.Err()))
			},
		}, {
			"single request, read next immediately, with error",
			func(t *testing.T, args testArgs) {
				args.mock.On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

				ch := args.ticker.Wait()

				_, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.False(t, ok)

				err, ok := readErrorFromChanWithTimeout(args.ticker.Err(), 30*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, "test error", err.Error())
			},
		},
	}

	// Test error handling.

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, mock := ethtesting.NewClientWithMock()

			ticker := ethhelpers.NewPeriodicBlockNumberTicker(ctx, client, 1, 100*time.Millisecond)
			defer ticker.Stop()

			test.fn(t, testArgs{
				ctx:    ctx,
				client: client,
				mock:   mock,
				ticker: ticker,
			})
		})
	}
}
