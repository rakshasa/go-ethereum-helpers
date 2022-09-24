package ethtesting

import (
	"github.com/ethereum/go-ethereum"
)

type ChainReaderAndTransactionSender interface {
	ethereum.TransactionSender
	ethereum.ChainReader
}
