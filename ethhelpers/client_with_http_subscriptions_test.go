package ethhelpers_test

import (
	"context"
	"fmt"
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

	client := ethhelpers.NewClientWithHTTPSubscriptions(ethtesting.NewSimulatedClient(sim.Backend))

	logChan := make(chan types.Log, 1)

	sub, err := client.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, logChan)
	if !assert.NoError(err) {
		return
	}

	commitAndNoErrorFn := func() bool {
		sim.Backend.Commit()
		time.Sleep(5 * time.Second)

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
		if !assert.NotEmpty(logChan) {
			return false
		}

		select {
		case log := <-logChan:
			return assert.Equal(txHash, log.TxHash)
		default:
			assert.Fail("log channel is empty")
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

	time.Sleep(5 * time.Second)

	if !assert.Empty(logChan) {
		fmt.Printf("got log: %+v", <-logChan)
	}

	select {
	case err, ok := <-sub.Err():
		assert.Nil(err)
		assert.False(ok)
	default:
		assert.Fail("error channel not closed")
	}
}
