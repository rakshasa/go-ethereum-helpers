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

func TestWaitForTransactionReceipt(t *testing.T) {
	assert := assert.New(t)

	commit := ethtesting.PendingLogHandlerForTesting(t, log.Root())
	defer commit()

	sim := ethtesting.NewSimulatedBackendWithAccounts(
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
		ErrorHandler: DefaultErrorHandlerWithMessages(func(msg string) {
			// TODO: Add a unit test for this.
			fmt.Printf("%s : %s\n", signedTx.Hash().String(), msg)
		}),
	})
	defer cancel()

	// TODO: Use custom time ticker.
	time.Sleep(5 * time.Second)
	assert.Empty(resultChan)

	sim.Backend.Commit()

	time.Sleep(5 * time.Second)
	assert.NotEmpty(resultChan)

	result := <-resultChan

	assert.Equal(uint64(1), result.Receipt.Status)
	assert.Equal(signedTx.Hash(), result.Receipt.TxHash)
	assert.Equal(common.Address{}, result.Receipt.ContractAddress)
	assert.Nil(result.Error)

	time.Sleep(5 * time.Second)
	assert.Empty(resultChan)

	// TODO: Test errors.
}
