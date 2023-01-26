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
	startBlock    uint64
	emptyResults  bool
	withFromBlock bool
	withTimestamp bool
	windowSize    uint64
}

// !!! test for big jump.
// stresstest
// use current time as fake timestamp for normal requests
// TODO: Make delta configurable.

func testTickers_periodic(t *testing.T, options tickersOptions) []blockNumberTickerTest {
	prefix := fmt.Sprintf("startBlock=%d", options.startBlock)

	if options.emptyResults {
		prefix += ",emptyResults"
	}
	if options.withTimestamp {
		prefix += ",withTimestamp"
	}

	callWait := func(ticker ethhelpers.BlockNumberTicker) interface{} {
		if !options.withTimestamp {
			t.Fatal("callWait called with withTimestamp=false")
			return nil
		} else {
			return ticker.Wait()
		}
	}

	emptyResult := func(t *testing.T, c interface{}) {
		if !options.withTimestamp {
			if ch, ok := c.(<-chan uint64); assert.True(t, ok) {
				select {
				case <-ch:
					t.Error("unexpected result")
				default:
				}
			}
		} else {
			if ch, ok := c.(<-chan ethhelpers.BlockNumber); assert.True(t, ok) {
				select {
				case <-ch:
					t.Error("unexpected result")
				default:
				}
			}
		}
	}
	emptyResultWithTimeout := func(t *testing.T, c interface{}, timeout time.Duration) {
		if !options.withTimestamp {
			if ch, ok := c.(<-chan uint64); assert.True(t, ok) {
				select {
				case <-ch:
					t.Error("unexpected result")
				case <-time.After(timeout):
				}
			}
		} else {
			if ch, ok := c.(<-chan ethhelpers.BlockNumber); assert.True(t, ok) {
				select {
				case <-ch:
					t.Error("unexpected result")
				case <-time.After(timeout):
				}
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
	expectedResult := func(t *testing.T, expectedBlock uint64, truncated bool, c interface{}) {
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
			if ch, ok := c.(<-chan ethhelpers.BlockNumber); assert.True(t, ok) {
				select {
				case r, ok := <-ch:
					assert.True(t, ok)
					assert.Equal(t, expectedBlock, r.BlockNumber)
					assert.WithinDuration(t, time.Now(), r.Timestamp, 300*time.Millisecond)
					assert.Equal(t, truncated, r.Truncated)
				default:
					assert.Fail(t, "expected result not received: channel empty")
				}
			}
		}
	}
	expectedResultWithTimeout := func(t *testing.T, expectedBlock uint64, truncated bool, c interface{}, timeout time.Duration) {
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
			if ch, ok := c.(<-chan ethhelpers.BlockNumber); assert.True(t, ok) {
				select {
				case r, ok := <-ch:
					assert.True(t, ok)
					assert.Equal(t, expectedBlock, r.BlockNumber)
					assert.WithinDuration(t, time.Now(), r.Timestamp, 300*time.Millisecond)
					assert.Equal(t, truncated, r.Truncated)
				case <-time.After(timeout):
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
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(200 * time.Millisecond)
			c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

			ch := callWait(ticker)
			time.Sleep(50 * time.Millisecond)

			emptyResult(t, ch)
			emptyError(t, ticker)
		},
		// }, blockNumberTickerTest{
		// "(" + prefix + ") handles context cancelation before wait call",
		// func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
		// 	c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

		// 	ch := callWait(ticker)
		// 	time.Sleep(50 * time.Millisecond)

		// 	emptyResult(t, ch)
		// 	emptyError(t, ticker)
	})

	// TODO: Add tests for context cancelation.

	if !options.emptyResults {
		//
		// Default Tests, with results:
		//

		tests = append(tests, blockNumberTickerTest{
			"(" + prefix + ") single request, read with timeout",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(20 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)
				emptyResult(t, ch)

				expectedResultWithTimeout(t, options.startBlock, false, ch, 50*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, read after sleeping",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(20 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)
				emptyResult(t, ch)

				time.Sleep(50 * time.Millisecond)

				expectedResult(t, options.startBlock, false, ch)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, read after sleeping beyond tick",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(20 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)

				time.Sleep(150 * time.Millisecond)

				expectedResult(t, options.startBlock, false, ch)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read next immediately, with new block number",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock+1), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock, false, ch, 30*time.Millisecond)

				ch = callWait(ticker)
				emptyResultWithTimeout(t, ch, 60*time.Millisecond)
				expectedResultWithTimeout(t, options.startBlock+1, false, ch, 70*time.Millisecond)

				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read after sleeping beyond tick, with new block number",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock+1), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock, false, ch, 30*time.Millisecond)

				time.Sleep(100 * time.Millisecond)

				ch = callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock+1, false, ch, 30*time.Millisecond)

				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read next immediately, with same block number",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Twice().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock, false, ch, 30*time.Millisecond)

				ch = callWait(ticker)
				emptyResultWithTimeout(t, ch, 150*time.Millisecond)

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
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).After(20 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)

				emptyResultWithTimeout(t, ch, 230*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, in same tick, empty results",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock+1), nil).Twice().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

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
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock+1), nil).Times(3).After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

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
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Twice().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

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
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				expectedErrorWithTimeout(t, ticker, 30*time.Millisecond)

				assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") double request, read after sleeping beyond tick with error, empty results",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

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
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock+1), nil).Once().After(10 * time.Millisecond)
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(10 * time.Millisecond)
			c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

			ch := callWait(ticker)

			if !options.emptyResults {
				expectedResultWithTimeout(t, options.startBlock+1, false, ch, 30*time.Millisecond)
			} else {
				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
			}

			time.Sleep(100 * time.Millisecond)

			ch = callWait(ticker)

			emptyResultWithTimeout(t, ch, 30*time.Millisecond)
			emptyError(t, ticker)
		},
	})

	//
	// Window size:
	//

	// TODO: Test with window size == 1 and 2.
	switch {
	case options.emptyResults:
		break

	case options.windowSize != 0:
		tests = append(tests, blockNumberTickerTest{
			"(" + prefix + ") single request, read four times, with 4x window size",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(options.startBlock, nil).Once().After(10 * time.Millisecond)

				ch := callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock, false, ch, 30*time.Millisecond)

				time.Sleep(100 * time.Millisecond)
				startBlock := options.startBlock + 1

				c.Mock().On("BlockNumber", mock.Anything).Return(startBlock+options.windowSize*uint64(4)-1, nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch = callWait(ticker)
				expectedResultWithTimeout(t, startBlock+options.windowSize-1, true, ch, 30*time.Millisecond)

				ch = callWait(ticker)
				expectedResultWithTimeout(t, startBlock+options.windowSize*2-1, true, ch, 10*time.Millisecond)

				ch = callWait(ticker)
				expectedResultWithTimeout(t, startBlock+options.windowSize*3-1, true, ch, 10*time.Millisecond)

				ch = callWait(ticker)
				expectedResultWithTimeout(t, startBlock+options.windowSize*4-1, false, ch, 10*time.Millisecond)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, exactly window size",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(options.startBlock, nil).Once().After(10 * time.Millisecond)

				ch := callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock, false, ch, 30*time.Millisecond)

				time.Sleep(100 * time.Millisecond)
				startBlock := options.startBlock + 1

				c.Mock().On("BlockNumber", mock.Anything).Return(startBlock+(options.windowSize-1), nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch = callWait(ticker)
				expectedResultWithTimeout(t, startBlock+(options.windowSize-1), false, ch, 30*time.Millisecond)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, exactly window size + 1",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(options.startBlock, nil).Once().After(10 * time.Millisecond)

				ch := callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock, false, ch, 30*time.Millisecond)

				time.Sleep(100 * time.Millisecond)
				startBlock := options.startBlock + 1

				c.Mock().On("BlockNumber", mock.Anything).Return(startBlock+(options.windowSize-1)+1, nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch = callWait(ticker)
				expectedResultWithTimeout(t, startBlock+(options.windowSize-1), true, ch, 30*time.Millisecond)

				ch = callWait(ticker)
				expectedResultWithTimeout(t, startBlock+(options.windowSize-1)+1, false, ch, 10*time.Millisecond)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)
			},
		}, blockNumberTickerTest{
			"(" + prefix + ") single request, exactly window size - 1",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(options.startBlock, nil).Once().After(10 * time.Millisecond)

				ch := callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock, false, ch, 30*time.Millisecond)

				time.Sleep(100 * time.Millisecond)
				startBlock := options.startBlock + 1

				c.Mock().On("BlockNumber", mock.Anything).Return(startBlock+(options.windowSize-1)-1, nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch = callWait(ticker)
				expectedResultWithTimeout(t, startBlock+(options.windowSize-1)-1, false, ch, 30*time.Millisecond)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)
			},
		})

		// TODO: Add test for a large jump.
		// TODO: Add tests with very large window size.

	default:
		tests = append(tests, blockNumberTickerTest{
			"(" + prefix + ") single request, read once over a large jump",
			func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
				c.Mock().On("BlockNumber", mock.Anything).Return(options.startBlock+1000000, nil).Once().After(10 * time.Millisecond)
				c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

				ch := callWait(ticker)
				expectedResultWithTimeout(t, options.startBlock+1000000, false, ch, 30*time.Millisecond)

				emptyResultWithTimeout(t, ch, 30*time.Millisecond)
				emptyError(t, ticker)
			},
		})
	}

	//
	// Errors:
	//

	tests = append(tests, blockNumberTickerTest{
		"(" + prefix + ") single request, read error immediately",
		func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)
			c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

			ch := callWait(ticker)

			expectedErrorWithTimeout(t, ticker, 30*time.Millisecond)
			emptyResult(t, ch)

			assert.True(t, closedErrorChanWithTimeout(ticker.Err(), 30*time.Millisecond))
		},
	}, blockNumberTickerTest{
		"(" + prefix + ") double request, read error after sleeping beyond tick",
		func(t *testing.T, c ethtesting.ClientWithMock, ticker ethhelpers.BlockNumberTicker) {
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(options.startBlock), nil).Once().After(10 * time.Millisecond)
			c.Mock().On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("test error")).Once().After(10 * time.Millisecond)
			c.Mock().On("BlockNumber", mock.Anything).Return(ethtesting.CanceledMockCall()).Maybe()

			ch := callWait(ticker)

			if !options.emptyResults {
				expectedResultWithTimeout(t, options.startBlock, false, ch, 30*time.Millisecond)
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
// TODO: Add 32-bit plus tests.

func TestTickers_NewPeriodicBlockNumberTicker(t *testing.T) {
	tests := []blockNumberTickerTest{}

	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    0,
		withTimestamp: true,
	})...)
	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    1,
		withTimestamp: true,
	})...)
	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    100000,
		withTimestamp: true,
	})...)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client := ethtesting.NewClientWithMock()
			client.Test(t)

			ticker := ethhelpers.NewPeriodicBlockNumberTicker(ctx, client, 100*time.Millisecond)
			defer ticker.Stop()

			test.fn(t, client, ticker)

			client.Mock().AssertExpectations(t)
		})
	}
}

func TestTickers_NewPeriodicBlockNumberTickerWithWindowSize_10(t *testing.T) {
	tests := []blockNumberTickerTest{}

	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    0,
		withTimestamp: true,
		windowSize:    10,
	})...)
	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    1,
		withTimestamp: true,
		windowSize:    10,
	})...)
	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    100000,
		withTimestamp: true,
		windowSize:    10,
	})...)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client := ethtesting.NewClientWithMock()
			client.Test(t)

			ticker := ethhelpers.NewPeriodicBlockNumberTickerWithWindowSize(ctx, client, 100*time.Millisecond, 10)
			defer ticker.Stop()

			test.fn(t, client, ticker)

			client.Mock().AssertExpectations(t)
		})
	}
}

func TestTickers_NewPeriodicBlockNumberTickerFromBlock_0(t *testing.T) {
	tests := []blockNumberTickerTest{}

	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    0,
		withFromBlock: true,
		withTimestamp: true,
	})...)
	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    1,
		withFromBlock: true,
		withTimestamp: true,
	})...)
	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    100000,
		withFromBlock: true,
		withTimestamp: true,
	})...)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client := ethtesting.NewClientWithMock()
			client.Test(t)

			ticker := ethhelpers.NewPeriodicBlockNumberTickerFromBlock(ctx, client, 100*time.Millisecond, 0)
			defer ticker.Stop()

			test.fn(t, client, ticker)

			client.Mock().AssertExpectations(t)
		})
	}
}

func TestTickers_NewPeriodicBlockNumberTickerFromBlock_100000(t *testing.T) {
	tests := []blockNumberTickerTest{}

	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    0,
		withFromBlock: true,
		withTimestamp: true,
		emptyResults:  true,
	})...)
	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    1,
		withFromBlock: true,
		withTimestamp: true,
		emptyResults:  true,
	})...)
	tests = append(tests, testTickers_periodic(t, tickersOptions{
		startBlock:    100000,
		withFromBlock: true,
		withTimestamp: true,
	})...)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client := ethtesting.NewClientWithMock()
			client.Test(t)

			ticker := ethhelpers.NewPeriodicBlockNumberTickerFromBlock(ctx, client, 100*time.Millisecond, 100000)
			defer ticker.Stop()

			test.fn(t, client, ticker)

			client.Mock().AssertExpectations(t)
		})
	}
}
