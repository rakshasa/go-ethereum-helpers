package ethhelpers_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/stretchr/testify/assert"
)

func TestWaitForTransactionReceipt(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	sim, closeSim := newDefaultSimulatedBackend(t)
	defer closeSim()

	ctx := context.Background()

	signedTx, err := sim.Accounts[0].SendNewTransaction(ctx, sim.Backend, sim.Accounts[0].NonceAndIncrement(), sim.Accounts[1].Address, big.NewInt(1000), params.TxGas, nil)
	if !assert.NoError(err) {
		return
	}

	resultChan, cancel := ethhelpers.WaitForTransactionReceipt(ctx, ethhelpers.WaitForTransactionReceiptOptions{
		// TODO: If Client is nil get from ctx.
		Client: sim.Backend,
		TxHash: signedTx.Hash(),
		ErrorHandler: ethhelpers.DefaultErrorHandlerWithMessages(func(txHash common.Hash, msg string) {
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
	t.Parallel()

	assert := assert.New(t)

	sim, closeSim := newDefaultSimulatedBackend(t)
	defer closeSim()

	ctx := context.Background()

	waiter := ethhelpers.NewWaitTransactionReceipts(ctx, func(ctx context.Context, txHash common.Hash) (<-chan ethhelpers.ReceiptOrError, func()) {
		return ethhelpers.WaitForTransactionReceipt(ctx, ethhelpers.WaitForTransactionReceiptOptions{
			Client: sim.Backend,
			TxHash: txHash,
			ErrorHandler: ethhelpers.DefaultErrorHandlerWithMessages(func(txHash common.Hash, msg string) {
				fmt.Printf("%s : %s\n", txHash, msg)

				// assert.Equal(signedTx.Hash(), txHash)
				assert.NotEmpty(msg)
			}),
		})
	})
	defer waiter.Stop()

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

	if signedTx1.Hash() == result.Receipt.TxHash {
		assert.Equal(signedTx1.Hash(), result.Receipt.TxHash)
	} else {
		assert.Equal(signedTx2.Hash(), result.Receipt.TxHash)
	}

	assert.Equal(common.Address{}, result.Receipt.ContractAddress)
	assert.Nil(result.Error)

	time.Sleep(5 * time.Second)
	assert.Empty(waiter.Result())
}
