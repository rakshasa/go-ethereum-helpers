package ethhelpers_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

func newTestDefaultSimulatedBackend(t *testing.T) (*ethtesting.SimulatedBackendWithAccounts, func()) {
	commitLogs := ethtesting.PendingLogHandlerForTesting(t, log.Root())

	sim := ethtesting.NewSimulatedBackendWithAccounts(
		ethtesting.GenesisAccountWithPrivateKey{
			PrivateKey: ethtesting.MockPrivateKey1,
			GenesisAccount: core.GenesisAccount{
				Balance: big.NewInt(10_000_000_000_000_000),
			},
		},
		ethtesting.GenesisAccountWithPrivateKey{
			PrivateKey: ethtesting.MockPrivateKey2,
			GenesisAccount: core.GenesisAccount{
				Balance: big.NewInt(10_000_000_000_000_000),
			},
		},
	)

	return sim, func() {
		sim.Backend.Close()
		commitLogs()
	}
}

func newTestDefaultSimulatedBackendWithCallableContract(t *testing.T) (*ethtesting.SimulatedBackendWithAccounts, *bind.BoundContract, func()) {
	sim, cancel := newTestDefaultSimulatedBackend(t)

	contract, err := sim.GenerateCallableContract()
	assert.NoError(t, err)

	return sim, contract, cancel
}

func readLogFromChan(ch <-chan types.Log) (types.Log, bool) {
	select {
	case r := <-ch:
		return r, true
	default:
		return types.Log{}, false
	}
}

func readUint64FromChan(ch <-chan uint64) (uint64, bool) {
	select {
	case r := <-ch:
		return r, true
	default:
		return 0, false
	}
}

func sendTestTransaction(ctx context.Context, sim *ethtesting.SimulatedBackendWithAccounts) (*types.Transaction, error) {
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
