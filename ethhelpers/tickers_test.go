package ethhelpers_test

import (
	"context"
	"testing"
	"time"

	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

func TestTickers_NewBlockNumberTickerWithDuration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// TODO: Add tester client that simulates errors.
	// TODO: Add test that wait with <1/4s sleep.
	// TODO: Test with repeated wait calls.
	// TODO: Test errors from client.BlockNumber(ctx)

	tests := []struct {
		name string
		fn   func(t *testing.T, sim *ethtesting.SimulatedBackendWithAccounts, ticker ethhelpers.BlockNumberTicker)
	}{
		{
			"new ticker",
			func(t *testing.T, sim *ethtesting.SimulatedBackendWithAccounts, ticker ethhelpers.BlockNumberTicker) {
				assert.Empty(t, ticker.Err())
				assert.Empty(t, ticker.Wait())

				ch := ticker.Wait()
				time.Sleep(time.Second)

				assert.Empty(t, ticker.Err())
				_, ok := readUint64FromChan(ch)
				assert.False(t, ok)
			},
		}, {
			"single block",
			func(t *testing.T, sim *ethtesting.SimulatedBackendWithAccounts, ticker ethhelpers.BlockNumberTicker) {
				sim.Backend.Commit()

				ch := ticker.Wait()
				time.Sleep(time.Second)

				assert.Empty(t, ticker.Err())
				r, ok := readUint64FromChan(ch)
				assert.True(t, ok)
				assert.Equal(t, uint64(1), r)

				ch = ticker.Wait()
				time.Sleep(time.Second)

				assert.Empty(t, ticker.Err())
				_, ok = readUint64FromChan(ch)
				assert.False(t, ok)
			},
		}, {
			"multiple blocks return latest",
			func(t *testing.T, sim *ethtesting.SimulatedBackendWithAccounts, ticker ethhelpers.BlockNumberTicker) {
				sim.Backend.Commit()
				time.Sleep(time.Second)
				sim.Backend.Commit()

				ch := ticker.Wait()
				time.Sleep(time.Second)

				assert.Empty(t, ticker.Err())
				r, ok := readUint64FromChan(ch)
				assert.True(t, ok)
				assert.Equal(t, uint64(2), r)
			},
		}, {
			"result before ticker triggers",
			func(t *testing.T, sim *ethtesting.SimulatedBackendWithAccounts, ticker ethhelpers.BlockNumberTicker) {
				sim.Backend.Commit()

				select {
				case <-ticker.Wait():
				case <-ctx.Done():
					assert.Fail(t, "timed out")
					return
				}

				sim.Backend.Commit()

				// TODO: Also test with two calls.
				ch := ticker.Wait()

				select {
				case <-ch:
					assert.Fail(t, "unexpected result before interval ticker")
				case <-time.After(time.Second / 8):
				}

				select {
				case <-ch:
				case <-time.After(time.Second / 4):
					assert.Fail(t, "no result after interval ticker")
				}
			},
		}, {
			"tick happens while having no result",
			func(t *testing.T, sim *ethtesting.SimulatedBackendWithAccounts, ticker ethhelpers.BlockNumberTicker) {
				sim.Backend.Commit()

				// TODO: Make this a test option.
				select {
				case <-ticker.Wait():
				case <-ctx.Done():
					assert.Fail(t, "timed out")
					return
				}

				ch := ticker.Wait()

				time.Sleep((3 * time.Second) / 8)

				sim.Backend.Commit()

				time.Sleep(time.Second / 4)

				select {
				case r := <-ch:
					assert.Equal(t, uint64(2), r)
				case <-time.After(time.Second / 4):
					assert.Fail(t, "no result after interval ticker")
				}
			},
		},

		// TODO: Request before ticker timeout.
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// TODO: Add non-zero start block numbers
			sim, closeSim := newTestDefaultSimulatedBackend(t)
			defer closeSim()

			ticker := ethhelpers.NewPeriodicBlockNumberTicker(ctx, ethtesting.NewSimulatedClient(sim.Backend), 1, time.Second/4)
			defer ticker.Stop()

			test.fn(t, sim, ticker)
		})
	}
}
