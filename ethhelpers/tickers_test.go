package ethhelpers_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type blockNumberTickerTest struct {
	name string
	fn   func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker)
}

type tickersOptions struct {
	fromBlock    uint64
	emptyResults bool
}

// TODO: Move prefix and fromBlock to options.

func testTickers_periodic_fromBlock(t *testing.T, options tickersOptions) []blockNumberTickerTest {
	// TODO: Add tester client that simulates errors.
	// TODO: Test with repeated wait calls.
	// TODO: Test fromBlock variants.

	// TODO: Test if starting before fromblock.

	prefix := fmt.Sprintf("fromBlock=%d", options.fromBlock)

	if options.emptyResults {
		prefix += ",emptyResults"
	}

	tests := []blockNumberTickerTest{}

	tests = append(tests, blockNumberTickerTest{
		"(" + prefix + ") has empty channels immediately after wait call",
		func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
			c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(200 * time.Millisecond)

			ch := ticker.Wait()
			time.Sleep(50 * time.Millisecond)

			assert.True(t, emptyUint64Channel(ch))
			assert.True(t, emptyErrorChannel(ticker.Err()))
		},
	})

	if !options.emptyResults {
		//
		// Default Tests, with results:
		//

		tests = append(tests, blockNumberTickerTest{
			"(" + prefix + ") single request, read with timeout",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(20 * time.Millisecond)

				ch := ticker.Wait()
				assert.True(t, emptyUint64Channel(ch))

				bn, ok := readUint64FromChanWithTimeout(ch, 50*time.Millisecond)
				_ = assert.True(t, ok) && assert.Equal(t, uint64(options.fromBlock), bn)

				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, read after sleeping",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(20 * time.Millisecond)

				ch := ticker.Wait()
				assert.True(t, emptyUint64Channel(ch))

				time.Sleep(50 * time.Millisecond)

				bn, ok := readUint64FromChan(ch)
				_ = assert.True(t, ok) && assert.Equal(t, uint64(options.fromBlock), bn)

				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, read after sleeping beyond tick",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(20 * time.Millisecond)

				ch := ticker.Wait()

				time.Sleep(150 * time.Millisecond)

				bn, ok := readUint64FromChan(ch)
				_ = assert.True(t, ok) && assert.Equal(t, uint64(options.fromBlock), bn)

				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read next immediately, with new block number",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Once().After(10 * time.Millisecond)

				ch := ticker.Wait()

				bn, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				_ = assert.True(t, ok) && assert.Equal(t, uint64(options.fromBlock), bn)

				ch = ticker.Wait()

				bn, ok = readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.False(t, ok)
				assert.Equal(t, uint64(0), bn)

				time.Sleep(100 * time.Millisecond)

				bn, ok = readUint64FromChan(ch)
				assert.True(t, ok)
				assert.Equal(t, uint64(options.fromBlock+1), bn)

				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read after sleeping beyond tick, with new block number",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Once().After(10 * time.Millisecond)

				ch := ticker.Wait()

				bn, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				_ = assert.True(t, ok) && assert.Equal(t, uint64(options.fromBlock), bn)

				time.Sleep(100 * time.Millisecond)

				ch = ticker.Wait()

				bn, ok = readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, uint64(options.fromBlock+1), bn)

				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read next immediately, with same block number",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Twice().After(10 * time.Millisecond)

				ch := ticker.Wait()

				bn, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				_ = assert.True(t, ok) && assert.Equal(t, uint64(options.fromBlock), bn)

				ch = ticker.Wait()

				_, ok = readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.False(t, ok)

				time.Sleep(100 * time.Millisecond)

				_, ok = readUint64FromChan(ch)
				assert.False(t, ok)

				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		})

	} else {
		//
		// Default Tests, no results:
		//

		tests = append(tests, blockNumberTickerTest{
			"(" + prefix + ") single request, wait for results",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Times(3).After(20 * time.Millisecond)

				ch := ticker.Wait()

				assert.True(t, emptyUint64ChannelWithTimeout(ch, 230*time.Millisecond))
				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, in same tick",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Twice().After(10 * time.Millisecond)

				ch := ticker.Wait()

				assert.True(t, emptyUint64ChannelWithTimeout(ch, 30*time.Millisecond))
				assert.True(t, emptyErrorChannel(ticker.Err()))

				ch = ticker.Wait()

				assert.True(t, emptyUint64ChannelWithTimeout(ch, 230*time.Millisecond))
				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, in next tick with next block number",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Times(3).After(10 * time.Millisecond)

				ch := ticker.Wait()

				assert.True(t, emptyUint64ChannelWithTimeout(ch, 30*time.Millisecond))
				assert.True(t, emptyErrorChannel(ticker.Err()))

				time.Sleep(100 * time.Millisecond)

				ch = ticker.Wait()

				assert.True(t, emptyUint64ChannelWithTimeout(ch, 230*time.Millisecond))
				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, in same tick with same block number",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Twice().After(10 * time.Millisecond)

				ch := ticker.Wait()

				assert.True(t, emptyUint64ChannelWithTimeout(ch, 30*time.Millisecond))
				assert.True(t, emptyErrorChannel(ticker.Err()))

				ch = ticker.Wait()

				assert.True(t, emptyUint64ChannelWithTimeout(ch, 230*time.Millisecond))
				assert.True(t, emptyErrorChannel(ticker.Err()))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, read next immediately with error",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

				ch := ticker.Wait()

				_, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.False(t, ok)

				err, ok := readErrorFromChanWithTimeout(ticker.Err(), 30*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, "test error", err.Error())

				assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read after sleeping beyond tick with error",
			func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.mock.On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

				ch := ticker.Wait()

				assert.True(t, emptyUint64ChannelWithTimeout(ch, 30*time.Millisecond))
				assert.True(t, emptyErrorChannel(ticker.Err()))

				time.Sleep(100 * time.Millisecond)

				ch = ticker.Wait()

				err, ok := readErrorFromChanWithTimeout(ticker.Err(), 30*time.Millisecond)
				assert.True(t, ok)
				assert.Equal(t, "test error", err.Error())

				assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
			},
		})
	}

	//
	// Errors:
	//

	tests = append(tests, blockNumberTickerTest{
		"(" + prefix + ") single request, read next immediately with error",
		func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
			c.mock.On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

			ch := ticker.Wait()

			_, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
			assert.False(t, ok)

			err, ok := readErrorFromChanWithTimeout(ticker.Err(), 30*time.Millisecond)
			_ = assert.True(t, ok) && assert.Equal(t, "test error", err.Error())

			assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
		},
	}, blockNumberTickerTest{
		"(" + prefix + ") double request, read after sleeping beyond tick with error",
		func(t *testing.T, c ClientAndMock, ticker ethhelpers.BlockNumberTicker) {
			c.mock.On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
			c.mock.On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

			ch := ticker.Wait()

			if !options.emptyResults {
				bn, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				_ = assert.True(t, ok) && assert.Equal(t, uint64(options.fromBlock), bn)
			} else {
				_, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
				assert.False(t, ok)
			}

			time.Sleep(100 * time.Millisecond)

			ch = ticker.Wait()

			err, ok := readErrorFromChanWithTimeout(ticker.Err(), 30*time.Millisecond)
			_ = assert.True(t, ok) && assert.Equal(t, "test error", err.Error())

			assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
		},
	})

	return tests
}

// TODO: Add FromBlock tests, add NewBlockNumberTickerWithFromBlockAndWindowSize?

func TestTickers_NewPeriodicBlockNumberTicker(t *testing.T) {
	tests := []blockNumberTickerTest{}

	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock: 0,
	})...)
	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock: 1,
	})...)
	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock: 100000,
	})...)

	// TODO: Add 32-bit plus tests.

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client, mock := ethtesting.NewClientWithMock()

			ticker := ethhelpers.NewPeriodicBlockNumberTicker(ctx, client, 100*time.Millisecond)
			defer ticker.Stop()

			test.fn(t, ClientAndMock{ctx, client, mock}, ticker)
		})
	}
}

func TestTickers_NewPeriodicBlockNumberTickerFromBlock_fromBlock_0(t *testing.T) {
	tests := []blockNumberTickerTest{}

	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock: 0,
	})...)
	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock: 1,
	})...)
	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock: 100000,
	})...)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client, mock := ethtesting.NewClientWithMock()

			ticker := ethhelpers.NewPeriodicBlockNumberTickerFromBlock(ctx, client, 100*time.Millisecond, 0)
			defer ticker.Stop()

			test.fn(t, ClientAndMock{ctx, client, mock}, ticker)
		})
	}
}

func TestTickers_NewPeriodicBlockNumberTickerFromBlock_fromBlock_100000(t *testing.T) {
	tests := []blockNumberTickerTest{}

	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock:    0,
		emptyResults: true,
	})...)
	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock:    1,
		emptyResults: true,
	})...)
	tests = append(tests, testTickers_periodic_fromBlock(t, tickersOptions{
		fromBlock: 100000,
	})...)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client, mock := ethtesting.NewClientWithMock()

			ticker := ethhelpers.NewPeriodicBlockNumberTickerFromBlock(ctx, client, 100*time.Millisecond, 100000)
			defer ticker.Stop()

			test.fn(t, ClientAndMock{ctx, client, mock}, ticker)
		})
	}
}
