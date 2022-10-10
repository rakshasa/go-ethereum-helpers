package ethtesting

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
)

func SimulatedChainID() *big.Int {
	return big.NewInt(1337)
}

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

const callableAbi = "[{\"anonymous\":false,\"inputs\":[],\"name\":\"Called\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

const callableBin = "6080604052348015600f57600080fd5b5060998061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806334e2292114602d575b600080fd5b60336035565b005b7f81fab7a4a0aa961db47eefc81f143a5220e8c8495260dd65b1356f1d19d3c7b860405160405180910390a156fea2646970667358221220029436d24f3ac598ceca41d4d712e13ced6d70727f4cdc580667de66d2f51d8b64736f6c63430008010033"

// GenerateCallableContract creates a callable contract on the simulated backend.
//
// The transaction must be commited by the caller.
//
//   pragma solidity >=0.7.0 <0.9.0;
//   contract Callable {
//       event Called();
//       function Call() public { emit Called(); }
//   }
//
func (b *SimulatedBackendWithAccounts) GenerateCallableContract() (*bind.BoundContract, error) {
	parsed, _ := abi.JSON(strings.NewReader(callableAbi))
	auth, _ := bind.NewKeyedTransactorWithChainID(MockPrivateKey1, SimulatedChainID())

	_, _, contract, err := bind.DeployContract(auth, parsed, common.FromHex(callableBin), b.Backend)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy callable contract: %v", err)
	}

	return contract, nil
}
