package ethhelpers_test

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"

	ethlog "github.com/ethereum/go-ethereum/log"
)

func newExampleDefaultSimulatedBackend() (*ethtesting.SimulatedBackendWithAccounts, func()) {
	oldHandler := ethlog.Root().GetHandler()
	ethlog.Root().SetHandler(ethlog.DiscardHandler())

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

	return sim, func() {
		sim.Backend.Close()
		ethlog.Root().SetHandler(oldHandler)
	}
}

func sendExampleTransaction(ctx context.Context, sim *ethtesting.SimulatedBackendWithAccounts) (*types.Transaction, error) {
	return sim.Accounts[0].SendNewTransaction(
		ctx,
		sim.Backend,
		sim.Accounts[0].NonceAndIncrement(),
		sim.Accounts[1].Address,
		big.NewInt(1000),
		params.TxGas,
		nil,
	)
}

func ExampleWaitForTransactionReceipt() {
	ctx := context.Background()

	sim, closeSim := newExampleDefaultSimulatedBackend()
	defer closeSim()

	signedTx, err := sendExampleTransaction(ctx, sim)
	if err != nil {
		log.Fatal(err)
	}

	resultChan, cancel := ethhelpers.WaitForTransactionReceipt(ctx, ethhelpers.WaitForTransactionReceiptOptions{
		Client: sim.Backend,
		TxHash: signedTx.Hash(),
		ErrorHandler: ethhelpers.DefaultErrorHandlerWithMessages(func(txHash common.Hash, msg string) {
		}),
	})
	defer cancel()

	sim.Backend.Commit()

	result := <-resultChan
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	_ = result.Receipt

	// Output:
}

func ExampleWaitTransactionReceipts() {
	ctx := context.Background()

	sim, closeSim := newExampleDefaultSimulatedBackend()
	defer closeSim()

	waiter := ethhelpers.NewWaitTransactionReceipts(ctx, func(ctx context.Context, txHash common.Hash) (<-chan ethhelpers.ReceiptOrError, func()) {
		return ethhelpers.WaitForTransactionReceipt(ctx, ethhelpers.WaitForTransactionReceiptOptions{
			Client: sim.Backend,
			TxHash: txHash,
			ErrorHandler: ethhelpers.DefaultErrorHandlerWithMessages(func(txHash common.Hash, msg string) {
			}),
		})
	})
	defer waiter.Stop()

	signedTx1, err := sendExampleTransaction(ctx, sim)
	if err != nil {
		log.Fatal(err)
	}
	waiter.Add(signedTx1.Hash())

	signedTx2, err := sendExampleTransaction(ctx, sim)
	if err != nil {
		log.Fatal(err)
	}
	waiter.Add(signedTx2.Hash())

	sim.Backend.Commit()

	result := <-waiter.Result()
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	_ = result.Receipt

	// waiter.Result() will not return any more results.

	// Output:
}
