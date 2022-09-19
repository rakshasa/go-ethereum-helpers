package ethhelpers

import (
	"math/big"
)

type Config struct {
	Endpoint  string
	ChainId   *big.Int
	Contracts ContractContainer
}
