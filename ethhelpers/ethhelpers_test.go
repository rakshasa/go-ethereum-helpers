package ethhelpers_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientAndMock struct {
	ctx    context.Context
	client ethtesting.ClientWithMock
	mock   *mock.Mock
}

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
	case r, ok := <-ch:
		return r, ok
	default:
		return types.Log{}, false
	}
}

func readUint64FromChan(ch <-chan uint64) (uint64, bool) {
	select {
	case r, ok := <-ch:
		return r, ok
	default:
		return 0, false
	}
}

func readUint64FromChanWithTimeout(ch <-chan uint64, after time.Duration) (uint64, bool) {
	select {
	case r, ok := <-ch:
		return r, ok
	case <-time.After(after):
		return 0, false
	}
}

func readErrorFromChan(ch <-chan error) (error, bool) {
	select {
	case r := <-ch:
		if r == nil {
			return nil, false
		}
		return r, true
	default:
		return nil, false
	}
}

func readErrorFromChanWithTimeout(ch <-chan error, after time.Duration) (error, bool) {
	select {
	case r := <-ch:
		if r == nil {
			return nil, false
		}
		return r, true
	case <-time.After(after):
		return nil, false
	}
}

func closedErrorChanWithTimeout(ch <-chan error, after time.Duration) bool {
	select {
	case _, ok := <-ch:
		return !ok
	case <-time.After(after):
		return false
	}
}

func emptyUint64Channel(ch <-chan uint64) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}

func emptyUint64ChannelWithTimeout(ch <-chan uint64, after time.Duration) bool {
	select {
	case <-ch:
		return false
	case <-time.After(after):
		return true
	}
}

func emptyErrorChannel(ch <-chan error) bool {
	select {
	case <-ch:
		return false
	default:
		return true
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

func unknownErrorWithAssertFail(t *testing.T) func(context.Context, error) error {
	return func(ctx context.Context, err error) error {
		assert.Failf(t, "should not be called", "%v", err)
		return context.Canceled
	}
}
