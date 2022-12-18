package ethhelpers_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

func TestClientWithHTTPSubscriptions_SubscribeFilterLogs(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	sim, contract, closeSim := newTestDefaultSimulatedBackendWithCallableContract(t)
	defer closeSim()

	// TODO: Speed up the test.
	client := ethhelpers.NewClientWithHTTPSubscriptions(
		ethtesting.NewSimulatedClient(sim.Backend),

		func(ctx context.Context, fromBlock uint64) ethhelpers.BlockNumberTicker {
			return ethhelpers.NewPeriodicBlockNumberTickerFromBlock(ctx, ethtesting.NewSimulatedClient(sim.Backend), time.Second/4, fromBlock)
		},
	)

	logChan := make(chan types.Log, 1)

	sub, err := client.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, logChan)
	if !assert.NoError(err) {
		return
	}

	time.Sleep(time.Second)

	commitAndNoErrorFn := func() bool {
		sim.Backend.Commit()
		time.Sleep(time.Second)

		select {
		case err, ok := <-sub.Err():
			assert.Nil(err)
			assert.True(ok)
			assert.Fail("error channel was not empty")
			return false
		default:
			return true
		}
	}
	verifyWithTxHashFn := func(txHash common.Hash) bool {
		select {
		case log, ok := <-logChan:
			return assert.True(ok) && assert.Equal(txHash, log.TxHash)
		case err, ok := <-sub.Err():
			assert.NoError(err)
			assert.False(ok)
			return false
		case <-ctx.Done():
			assert.Fail("timed out")
			return false
		}
	}

	if !commitAndNoErrorFn() || !assert.Empty(logChan) {
		return
	}

	signer1, _ := bind.NewKeyedTransactorWithChainID(ethtesting.MockPrivateKey1, ethtesting.SimulatedChainID())

	tx1, err := contract.Transact(signer1, "Call")
	if !assert.NoError(err) {
		return
	}
	if !commitAndNoErrorFn() || !verifyWithTxHashFn(tx1.Hash()) {
		return
	}

	signer2, _ := bind.NewKeyedTransactorWithChainID(ethtesting.MockPrivateKey2, ethtesting.SimulatedChainID())

	tx2, err := contract.Transact(signer2, "Call")
	if !assert.NoError(err) {
		return
	}
	if !commitAndNoErrorFn() || !verifyWithTxHashFn(tx2.Hash()) {
		return
	}

	sub.Unsubscribe()

	time.Sleep(time.Second)

	if _, ok := readLogFromChan(logChan); !assert.False(ok) {
		return
	}

	select {
	case err := <-sub.Err():
		assert.Error(err)
	case <-ctx.Done():
		assert.Fail("timed out")
	}

	select {
	case _, ok := <-sub.Err():
		assert.False(ok)
	case <-ctx.Done():
		assert.Fail("timed out")
	}
}
