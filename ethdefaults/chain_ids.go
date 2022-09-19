package ethdefaults

import (
	"math/big"
)

var (
	chainIdForEthereumMainnet          = big.NewInt(1)
	chainIdForBinanceSmartChainMainnet = big.NewInt(56)
	chainIdForPolygonMainnet           = big.NewInt(137)
	chainIdForArbitrumOne              = big.NewInt(42161)
)

func ChainIdForEthereumMainnet() *big.Int {
	return new(big.Int).Set(chainIdForEthereumMainnet)
}

func ChainIdForBinanceSmartChainMainnet() *big.Int {
	return new(big.Int).Set(chainIdForBinanceSmartChainMainnet)
}

func ChainIdForPolygonMainnet() *big.Int {
	return new(big.Int).Set(chainIdForPolygonMainnet)
}

func ChainIdForArbitrumOne() *big.Int {
	return new(big.Int).Set(chainIdForArbitrumOne)
}
