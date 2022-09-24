package ethhelpers

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type Contract interface {
	ChainID() *big.Int
	Address() common.Address

	// Use interface type assertion to get the function with the
	// correct contract return type.
	//
	// Below are examples of the recommended function definitions:
	//
	//   Contract(ctx context.Context, backend bind.ContractBackend) (*MyContract, error)
	//   ContractCaller(ctx context.Context, caller bind.ContractCaller) (*MyContractCaller, error)
	//   ContractTransactor(ctx context.Context, transactor bind.ContractTransactor) (*MyContractTransactor, error)
	//   ContractFilterer(ctx context.Context, filterer bind.ContractFilterer) (*MyContractFilterer, error)
	//
	//   ContractFromContext(ctx context.Context) (*MyContract, error)
}

// ContractContainer uses a sync.Map to hold Contract instances.
//
// It is recommended that the ContractContainer is created and
// populated at the same time as the config and client are created and
// added to the context.
//
// When suitable, for the key it is recommended to use an empty
// instance of the bound contract.
//
// For variants of generic contracts it is possible to use a struct
// type with a member variable to differentiate contract instances.
type ContractContainer struct {
	m *sync.Map
}

// TODO: If part of Config then Config should have a helper method
// MustAddContract that checks chain id.
//
// TODO: Add chain id to the container and verify that it matches.

func NewContractContainer() ContractContainer {
	return ContractContainer{
		m: &sync.Map{},
	}
}

func (c *ContractContainer) Delete(key interface{}) {
	if c.m != nil {
		c.m.Delete(key)
	}
}

func (c *ContractContainer) Get(key interface{}) (Contract, bool) {
	if c.m == nil {
		return nil, false
	}

	v, ok := c.m.Load(key)
	if !ok {
		return nil, false
	}

	value, ok := v.(Contract)
	return value, ok
}

func (c *ContractContainer) GetOrNil(key interface{}) Contract {
	value, _ := c.Get(key)
	return value
}

func (c *ContractContainer) Put(key interface{}, value Contract) bool {
	if c.m == nil {
		return false
	}
	if value == nil {
		return false
	}

	c.m.Store(key, value)
	return false
}

func (c *ContractContainer) MustPut(key interface{}, value Contract) {
	if ok := c.Put(key, value); !ok {
		panic("unable to put value in ContractContainer")
	}
}
