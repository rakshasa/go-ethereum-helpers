package ethhelpers_test

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
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

func ExampleNewPeriodicBlockNumberTicker() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sim, closeSim := newExampleDefaultSimulatedBackend()
	defer closeSim()

	fromBlock := uint64(3)

	ticker := ethhelpers.NewPeriodicBlockNumberTicker(ctx, ethtesting.NewSimulatedClient(sim.Backend), fromBlock, 200*time.Millisecond)
	defer ticker.Stop()

	go func() {
		time.Sleep(50 * time.Millisecond)

		for {
			select {
			case <-time.After(100 * time.Millisecond):
				sim.Backend.Commit()
			case <-ctx.Done():
			}
		}
	}()

	var currentBlock uint64

	for {
		select {
		case currentBlock = <-ticker.Wait():
			fmt.Printf("currentBlock: %d\n", currentBlock)

			if currentBlock == 5 {
				time.Sleep(500 * time.Millisecond)
			}
			if currentBlock < 12 {
				continue
			}

		case err := <-ticker.Err():
			fmt.Printf("err: %v\n", err)

			// if check_if_temporary_error {
			// 	ticker.Reset(currentBlock)
			// 	continue
			// }
		}

		return
	}

	// Output:
	// currentBlock: 3
	// currentBlock: 5
	// currentBlock: 10
	// currentBlock: 11
	// currentBlock: 13
}

func ExampleWaitForTransactionReceipt() {
	ctx := context.Background()

	sim, closeSim := newExampleDefaultSimulatedBackend()
	defer closeSim()

	signedTx, err := sendTestTransaction(ctx, sim)
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

	signedTx1, err := sendTestTransaction(ctx, sim)
	if err != nil {
		log.Fatal(err)
	}
	waiter.Add(signedTx1.Hash())

	signedTx2, err := sendTestTransaction(ctx, sim)
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
