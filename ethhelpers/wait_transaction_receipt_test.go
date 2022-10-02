package ethhelpers

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

func newTestDefaultSimulatedBackend() *ethtesting.SimulatedBackendWithAccounts {
	return ethtesting.NewSimulatedBackendWithAccounts(
		ethtesting.GenesisAccountWithPrivateKey{
			PrivateKey: ethtesting.MockPrivateKey1,
			GenesisAccount: core.GenesisAccount{
				Balance: big.NewInt(10_000_000_000_000_000),
			},
		},
		ethtesting.GenesisAccountWithPrivateKey{
			PrivateKey: ethtesting.MockPrivateKey2,
		},
	)
}

func TestWaitForTransactionReceipt(t *testing.T) {
	assert := assert.New(t)

	commitLogs := ethtesting.PendingLogHandlerForTesting(t, log.Root())
	defer commitLogs()

	sim := newTestDefaultSimulatedBackend()
	defer sim.Backend.Close()

	ctx := context.Background()

	signedTx, err := sim.Accounts[0].SendNewTransaction(ctx, sim.Backend, sim.Accounts[0].NonceAndIncrement(), sim.Accounts[1].Address, big.NewInt(1000), params.TxGas, nil)
	if !assert.NoError(err) {
		return
	}

	resultChan, cancel := WaitForTransactionReceipt(ctx, WaitForTransactionReceiptOptions{
		// TODO: If Client is nil get from ctx.
		Client: sim.Backend,
		TxHash: signedTx.Hash(),
		ErrorHandler: DefaultErrorHandlerWithMessages(func(txHash common.Hash, msg string) {
			// TODO: Add a unit test for this.
			fmt.Printf("%s : %s\n", txHash, msg)

			assert.Equal(signedTx.Hash(), txHash)
			assert.NotEmpty(msg)
		}),
	})
	defer cancel()

	// TODO: Use custom time ticker.
	time.Sleep(5 * time.Second)
	assert.Empty(resultChan)

	sim.Backend.Commit()

	time.Sleep(5 * time.Second)

	if !assert.NotEmpty(resultChan) {
		return
	}
	result := <-resultChan

	assert.Equal(uint64(1), result.Receipt.Status)
	assert.Equal(signedTx.Hash(), result.Receipt.TxHash)
	assert.Equal(common.Address{}, result.Receipt.ContractAddress)
	assert.Nil(result.Error)

	time.Sleep(5 * time.Second)
	assert.Empty(resultChan)
}

// TODO: Add a tester that uses a mock wait funcion.
func TestWaitTransactionReceipts(t *testing.T) {
	assert := assert.New(t)

	commitLogs := ethtesting.PendingLogHandlerForTesting(t, log.Root())
	defer commitLogs()

	sim := newTestDefaultSimulatedBackend()
	defer sim.Backend.Close()

	ctx := context.Background()

	waiter := NewWaitTransactionReceipts(ctx, func(ctx context.Context, txHash common.Hash) (<-chan ReceiptOrError, func()) {
		return WaitForTransactionReceipt(ctx, WaitForTransactionReceiptOptions{
			Client: sim.Backend,
			TxHash: txHash,
			ErrorHandler: DefaultErrorHandlerWithMessages(func(txHash common.Hash, msg string) {
				fmt.Printf("%s : %s\n", txHash, msg)

				// assert.Equal(signedTx.Hash(), txHash)
				assert.NotEmpty(msg)
			}),
		})
	})

	time.Sleep(5 * time.Second)
	assert.Empty(waiter.Result())

	signedTx1, err := sim.Accounts[0].SendNewTransaction(ctx, sim.Backend, sim.Accounts[0].NonceAndIncrement(), sim.Accounts[1].Address, big.NewInt(1000), params.TxGas, nil)
	if !assert.NoError(err) {
		return
	}
	waiter.Add(signedTx1.Hash())

	// TODO: Find a way to test different transactions succeeding.
	signedTx2, err := sim.Accounts[0].SendNewTransaction(ctx, sim.Backend, sim.Accounts[0].NonceAndIncrement(), sim.Accounts[1].Address, big.NewInt(1000), params.TxGas, nil)
	if !assert.NoError(err) {
		return
	}
	waiter.Add(signedTx2.Hash())

	time.Sleep(5 * time.Second)
	assert.Empty(waiter.Result())

	sim.Backend.Commit()

	time.Sleep(5 * time.Second)

	if !assert.NotEmpty(waiter.Result()) {
		return
	}
	result := <-waiter.Result()

	assert.Equal(uint64(1), result.Receipt.Status)
	assert.Equal(signedTx1.Hash(), result.Receipt.TxHash)
	assert.Equal(common.Address{}, result.Receipt.ContractAddress)
	assert.Nil(result.Error)

	time.Sleep(5 * time.Second)
	assert.Empty(waiter.Result())
}
