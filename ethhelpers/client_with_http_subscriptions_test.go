package ethhelpers_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

func TestClientWithHTTPSubscriptions_SubscribeFilterLogs(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sim, closeSim := newTestDefaultSimulatedBackend(t)
	defer closeSim()

	client := ethhelpers.NewClientWithHTTPSubscriptions(ethtesting.NewSimulatedClient(sim.Backend))

	// TODO: Start subscribing...

	ch := make(chan types.Log)

	_, err := client.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, ch)
	if !assert.NoError(err) {
		return
	}
}
