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
	fn   func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker)
}

type tickersOptions struct {
	fromBlock     uint64
	emptyResults  bool
	withTimestamp bool
	windowSize    uint64
}

// !!! test for big jump.
// stresstest

func testTickers_periodic_fromBlock(t *testing.T, options tickersOptions) []blockNumberTickerTest {
	prefix := fmt.Sprintf("fromBlock=%d", options.fromBlock)

	if options.emptyResults {
		prefix += ",emptyResults"
	}

	callWait := func(ticker ethhelpers.BlockNumberTicker) interface{} {
		if !options.withTimestamp {
			return ticker.Wait()
		} else {
			return ticker.WaitWithTimestamp()
		}
	}

	emptyResult := func(t *testing.T, c interface{}) {
		if ch, ok := c.(<-chan uint64); assert.True(t, ok) {
			select {
			case <-ch:
				t.Error("unexpected result")
			default:
			}
		}
	}
	emptyResultWithTimeout := func(t *testing.T, c interface{}, timeout time.Duration) {
		if ch, ok := c.(<-chan uint64); assert.True(t, ok) {
			select {
			case <-ch:
				t.Error("unexpected result")
			case <-time.After(timeout):
			}
		}
	}
	emptyError := func(t *testing.T, ticker ethhelpers.BlockNumberTicker) {
		select {
		case err := <-ticker.Err():
			t.Error("unexpected error", err)
		default:
		}
	}
	expectedResult := func(t *testing.T, expectedBlock uint64, c interface{}) {
		if !options.withTimestamp {
			if ch, ok := c.(<-chan uint64); assert.True(t, ok) {
				select {
				case r, ok := <-ch:
					assert.True(t, ok)
					assert.Equal(t, expectedBlock, r)
				default:
					assert.Fail(t, "expected result not received: channel empty")
				}
			}
		} else {
			if ch, ok := c.(<-chan ethhelpers.BlockNumberWithTimestamp); assert.True(t, ok) {
				select {
				case r, ok := <-ch:
					assert.True(t, ok)
					assert.Equal(t, expectedBlock, r.BlockNumber)
					// TODO: Make delta configurable.
					assert.WithinDuration(t, time.Now(), r.Timestamp, 300*time.Millisecond)
				default:
					assert.Fail(t, "expected result not received: channel empty")
				}
			}
		}
	}
	expectedResultWithTimeout := func(t *testing.T, expectedBlock uint64, c interface{}, timeout time.Duration) {
		if !options.withTimestamp {
			if ch, ok := c.(<-chan uint64); assert.True(t, ok) {
				select {
				case r, ok := <-ch:
					assert.True(t, ok)
					assert.Equal(t, expectedBlock, r)
				case <-time.After(timeout):
					assert.Fail(t, "expected result not received: timeout")
				}
			}
		} else {
			if ch, ok := c.(<-chan ethhelpers.BlockNumberWithTimestamp); assert.True(t, ok) {
				select {
				case r, ok := <-ch:
					assert.True(t, ok)
					assert.Equal(t, expectedBlock, r.BlockNumber)
					assert.WithinDuration(t, time.Now(), r.Timestamp, 300*time.Millisecond)
				default:
					assert.Fail(t, "expected result not received: timeout")
				}
			}
		}
	}
	expectedErrorWithTimeout := func(t *testing.T, ticker ethhelpers.BlockNumberTicker, timeout time.Duration) {
		select {
		case err, ok := <-ticker.Err():
			assert.True(t, ok)
			assert.Error(t, err)
		case <-time.After(timeout):
			assert.Fail(t, "expected error not received: timeout")
		}
	}

	tests := []blockNumberTickerTest{}

	tests = append(tests, blockNumberTickerTest{
		"(" + prefix + ") has empty channels immediately after wait call",
		func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(200 * time.Millisecond)

			ch := callWait(ticker)
			time.Sleep(50 * time.Millisecond)

			emptyResult(t, ch)
			emptyError(t, ticker)
		},
	})

	if !options.emptyResults {
		//
		// Default Tests, with results:
		//

		tests = append(tests, blockNumberTickerTest{
			"(" + prefix + ") single request, read with timeout",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(20 * time.Millisecond)

				ch := callWait(ticker)
				emptyResult(t, ch)

				expectedResultWithTimeout(t, uint64(options.fromBlock), ch, 50*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, read after sleeping",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(20 * time.Millisecond)

				ch := callWait(ticker)
				emptyResult(t, ch)

				time.Sleep(50 * time.Millisecond)

				expectedResult(t, uint64(options.fromBlock), ch)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, read after sleeping beyond tick",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(20 * time.Millisecond)

				ch := callWait(ticker)

				time.Sleep(150 * time.Millisecond)

				expectedResult(t, uint64(options.fromBlock), ch)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read next immediately, with new block number",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Once().After(10 * time.Millisecond)

				ch := callWait(ticker)
				expectedResultWithTimeout(t, uint64(options.fromBlock), ch, 30*time.Millisecond)

				ch = callWait(ticker)
				expectedResultWithTimeout(t, uint64(0), ch, 30*time.Millisecond)

				time.Sleep(100 * time.Millisecond)

				expectedResult(t, uint64(options.fromBlock+1), ch)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read after sleeping beyond tick, with new block number",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Once().After(10 * time.Millisecond)

				ch := callWait(ticker)
				expectedResultWithTimeout(t, uint64(options.fromBlock), ch, 30*time.Millisecond)

				time.Sleep(100 * time.Millisecond)

				ch = callWait(ticker)
				expectedResultWithTimeout(t, uint64(options.fromBlock+1), ch, 30*time.Millisecond)

				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read next immediately, with same block number",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Twice().After(10 * time.Millisecond)

				ch := callWait(ticker)
				expectedResultWithTimeout(t, uint64(options.fromBlock), ch, 30*time.Millisecond)

				ch = callWait(ticker)
				emptyResultWithTimeout(t, ch, 130*time.Millisecond)

				emptyError(t, ticker)
			},
		})

	} else {
		//
		// Default Tests, no results:
		//

		tests = append(tests, blockNumberTickerTest{
			"(" + prefix + ") single request, wait for empty results",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Times(3).After(20 * time.Millisecond)

				ch := callWait(ticker)

				emptyResultWithTimeout(t, ch, 230*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, in same tick, empty results",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Twice().After(10 * time.Millisecond)

				ch := callWait(ticker)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)

				ch = callWait(ticker)

				emptyResultWithTimeout(t, ch, 230*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, in next tick with next block number, empty results",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Times(3).After(10 * time.Millisecond)

				ch := callWait(ticker)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)

				time.Sleep(100 * time.Millisecond)

				ch = callWait(ticker)

				emptyResultWithTimeout(t, ch, 230*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, in same tick with same block number, empty results",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Twice().After(10 * time.Millisecond)

				ch := callWait(ticker)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)

				ch = callWait(ticker)

				emptyResultWithTimeout(t, ch, 230*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, read next immediately with error, empty results",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

				ch := callWait(ticker)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				expectedErrorWithTimeout(t, ticker, 30*time.Millisecond)

				assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read after sleeping beyond tick with error, empty results",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

				ch := callWait(ticker)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)

				time.Sleep(100 * time.Millisecond)

				ch = callWait(ticker)
				expectedErrorWithTimeout(t, ticker, 30*time.Millisecond)

				assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
			},
		})
	}

	//
	// Invalid states:
	//

	tests = append(tests, blockNumberTickerTest{
		"(" + prefix + ") double request, read after sleeping beyond tick, with backwards jump",
		func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+1), nil).Once().After(10 * time.Millisecond)
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)

			ch := callWait(ticker)

			expectedResultWithTimeout(t, uint64(options.fromBlock+1), ch, 30*time.Millisecond)
			time.Sleep(100 * time.Millisecond)

			ch = callWait(ticker)

			emptyResultWithTimeout(t, ch, 30*time.Millisecond)
			emptyError(t, ticker)
		},
	})

	//
	// Window size:
	//

	// TODO: Add tests with window size == 1

	// if options.windowSize != 0 && !options.emptyResults {
	// 	tests = append(tests, blockNumberTickerTest{
	// 		"(" + prefix + ") double request, read next immediately, past window size",
	// 		func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
	// 			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+options.windowSize), nil).Once().After(10 * time.Millisecond)
	// 			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock+options.windowSize+1), nil).Once().After(10 * time.Millisecond)

	// 			ch := callWait(ticker)

	// 			bn, ok := readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
	// 			assert.True(t, ok)
	//assert.Equal(t, uint64(options.fromBlock+options.windowSize), bn)

	// 			ch = callWait(ticker)

	// 			bn, ok = readUint64FromChanWithTimeout(ch, 30*time.Millisecond)
	// 			assert.True(t, ok)
	//assert.Equal(t, uint64(options.fromBlock+options.windowSize+1), bn)

	// 			emptyError(t, ticker)
	// 		},
	// 	})
	// }

	//
	// Errors:
	//

	tests = append(tests, blockNumberTickerTest{
		"(" + prefix + ") single request, read error immediately",
		func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

			ch := callWait(ticker)

			expectedErrorWithTimeout(t, ticker, 30*time.Millisecond)
			emptyResult(t, ch)

			assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
		},
	}, blockNumberTickerTest{
		"(" + prefix + ") double request, read error after sleeping beyond tick",
		func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.fromBlock), nil).Once().After(10 * time.Millisecond)
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)

			ch := callWait(ticker)

			if !options.emptyResults {
				expectedResultWithTimeout(t, uint64(options.fromBlock), ch, 30*time.Millisecond)
			} else {
				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
			}

			time.Sleep(100 * time.Millisecond)

			ch = callWait(ticker)

			expectedErrorWithTimeout(t, ticker, 30*time.Millisecond)

			assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
		},
	})

	return tests
}

// TODO: Add FromBlock tests, add NewBlockNumberTickerWithFromBlockAndWindowSize?
// TODO: Test Clone.

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

			client := ethtesting.NewClientWithMock()

			ticker := ethhelpers.NewPeriodicBlockNumberTicker(ctx, client, 100*time.Millisecond)
			defer ticker.Stop()

			test.fn(t, client, ticker)

			client.Mock().AssertExpectations(t)
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

			client := ethtesting.NewClientWithMock()

			ticker := ethhelpers.NewPeriodicBlockNumberTickerFromBlock(ctx, client, 100*time.Millisecond, 0)
			defer ticker.Stop()

			test.fn(t, client, ticker)

			client.Mock().AssertExpectations(t)
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

			client := ethtesting.NewClientWithMock()

			ticker := ethhelpers.NewPeriodicBlockNumberTickerFromBlock(ctx, client, 100*time.Millisecond, 100000)
			defer ticker.Stop()

			test.fn(t, client, ticker)

			client.Mock().AssertExpectations(t)
		})
	}
}
