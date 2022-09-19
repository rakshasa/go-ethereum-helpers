package ethhelpers

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type Contract interface {
	ChainId() *big.Int
	Address() common.Address

	// Use interface type assertion to get the function with the
	// correct contract return type.
	//
	// Below are examples of the recommended function definitions:
	//
	// TODO: Add ContractFromContext variant?
	//
	// Contract(backend bind.ContractBackend) (*MyContract, error)
	// ContractCaller(caller bind.ContractCaller) (*MyContractCaller, error)
	// ContractTransactor(transactor bind.ContractTransactor) (*MyContractTransactor, error)
	// ContractFilterer(filterer bind.ContractFilterer) (*MyContractFilterer, error)
}

// ContractContainer is uses a sync.Map to hold Contract instances for
// use with e.g. ContractsFromContext.
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
//
// TODO: If part of Config then Config should have a helper method
// MustAddContract that checks chain id.
type ContractContainer struct {
	m *sync.Map
}

func NewContractContainer() *ContractContainer {
	return &ContractContainer{
		m: &sync.Map{},
	}
}

func (c *ContractContainer) Delete(key interface{}) {
	c.m.Delete(key)
}

func (c *ContractContainer) Get(key interface{}) (Contract, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return nil, false
	}

	value, ok := v.(Contract)
	return value, ok
}

func (c *ContractContainer) Put(key interface{}, value Contract) {
	if value == nil {
		c.m.Delete(key)
		return
	}

	c.m.Store(key, value)
}
