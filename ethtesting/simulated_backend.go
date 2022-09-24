package ethtesting

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
)

type SimulatedBackendWithAccounts struct {
	Backend  *backends.SimulatedBackend
	Accounts []*SimulatedAccount
}

type GenesisAccountWithPrivateKey struct {
	*ecdsa.PrivateKey
	core.GenesisAccount
}

// NewSimulatedBackendWithAccounts creates a new simulated backend
// with provided simulated accounts.
//
// If GenesisAccount.Balance is nil then big.NewInt(0) is used.
func NewSimulatedBackendWithAccounts(genesisAccounts ...GenesisAccountWithPrivateKey) *SimulatedBackendWithAccounts {
	genesisAlloc := core.GenesisAlloc{}
	accounts := make([]*SimulatedAccount, len(genesisAccounts))

	for idx, a := range genesisAccounts {
		account := NewSimulatedAccountWithNonce(a.PrivateKey, a.GenesisAccount.Nonce)
		accounts[idx] = account

		if a.GenesisAccount.Balance == nil {
			a.GenesisAccount.Balance = big.NewInt(0)
		}

		genesisAlloc[account.Address] = a.GenesisAccount
	}

	return &SimulatedBackendWithAccounts{
		Backend:  backends.NewSimulatedBackend(genesisAlloc, 10000000),
		Accounts: accounts,
	}
}
