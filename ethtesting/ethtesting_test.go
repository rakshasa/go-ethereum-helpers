package ethtesting_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
)

func newDefaultSimulatedBackend(t *testing.T) (*ethtesting.SimulatedBackendWithAccounts, func()) {
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
