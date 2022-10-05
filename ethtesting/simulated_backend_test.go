package ethtesting_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

func TestSimulatedBackendWithAccounts(t *testing.T) {
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
	if !assert.NoError(t, err) {
		return
	}
	_, err = sim.Backend.TransactionReceipt(ctx, signedTx.Hash())
	if !assert.Error(t, err) {
		return
	}

	sim.Backend.Commit()

	receipt, err := sim.Backend.TransactionReceipt(ctx, signedTx.Hash())
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, uint64(1), receipt.Status)
	assert.Equal(t, signedTx.Hash(), receipt.TxHash)
	assert.Equal(t, common.Address{}, receipt.ContractAddress)
}
