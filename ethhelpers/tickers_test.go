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
	assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sim, closeSim := newTestDefaultSimulatedBackend(t)
	defer closeSim()

	ticker := ethhelpers.NewBlockNumberTickerWithDuration(ctx, ethtesting.NewSimulatedClient(sim.Backend), time.Second/4)
	defer ticker.Stop()

	assert.Empty(ticker.Err())
	assert.Empty(ticker.Wait())

	time.Sleep(2 * time.Second)

	assert.Empty(ticker.Err())
	assert.Empty(ticker.Wait())

	sim.Backend.Commit()

	time.Sleep(2 * time.Second)

	assert.Empty(ticker.Err())
	if !assert.NotEmpty(ticker.Wait()) {
		return
	}
	assert.Equal(uint64(1), <-ticker.Wait())

	time.Sleep(2 * time.Second)

	assert.Empty(ticker.Err())
	assert.Empty(ticker.Wait())

	sim.Backend.Commit()
	sim.Backend.Commit()

	time.Sleep(2 * time.Second)

	assert.Empty(ticker.Err())
	if !assert.NotEmpty(ticker.Wait()) {
		return
	}
	assert.Equal(uint64(3), <-ticker.Wait())
}
