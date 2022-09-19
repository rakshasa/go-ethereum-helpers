package ethtesting

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	DefaultBlockGeneratorOffsetTime = int64(5)
)

type BlockGenerator interface {
	Generator() (func(int, *core.BlockGen), int, error)
}

type TransactionWithHeight struct {
	Height      int
	Transaction *types.Transaction
}

type SimpleBlockGenerator struct {
	// Expected block height of the chain if it contains a non-zero
	// value, otherwise the last block with a transaction is
	// used.
	//
	// Results in an error if set to a non-zero value that is less
	// than the last block number with a transaction.
	ExpectedHeight int

	// Offset time of the block in seconds.
	OffsetTime int64

	// Generate a byte array for the extraData generator parameter field.
	ExtraGenerator func(int, *core.BlockGen) []byte

	Transactions []TransactionWithHeight
}

func (g *SimpleBlockGenerator) Generator() (func(int, *core.BlockGen), int, error) {
	offsetTime := g.OffsetTime
	if offsetTime == 0 {
		offsetTime = DefaultBlockGeneratorOffsetTime
	}

	extraGenerator := g.ExtraGenerator

	transactions, lastBlockHeightWithTransaction, err := g.verifyTransactions()
	if err != nil {
		return nil, 0, fmt.Errorf("invalid transaction order")
	}

	height := g.ExpectedHeight

	if height == 0 {
		height = lastBlockHeightWithTransaction
	}
	if height < lastBlockHeightWithTransaction {
		return nil, 0, fmt.Errorf("transactions exceed expected block height")
	}

	transactionIndex := 0

	// The generator closure should use copies of the parameters from SimpleBlockGenerator.
	return func(blockIndex int, blockGen *core.BlockGen) {
		blockGen.OffsetTime(offsetTime)

		if extraGenerator != nil {
			blockGen.SetExtra(extraGenerator(blockIndex, blockGen))
		}

		// Catch block index out-of-order calls.
		if transactionIndex != 0 && transactions[transactionIndex-1].Height >= blockIndex {
			transactionIndex = 0
		}

		for transactionIndex != len(transactions) && transactions[transactionIndex].Height <= blockIndex {
			if blockIndex == transactions[transactionIndex].Height {
				blockGen.AddTx(transactions[transactionIndex].Transaction)
			}

			transactionIndex++
		}
	}, height, nil
}

func (g *SimpleBlockGenerator) verifyTransactions() ([]TransactionWithHeight, int, error) {
	if g.Transactions == nil {
		return []TransactionWithHeight{}, 0, nil
	}

	height := 0

	for _, transaction := range g.Transactions {
		if transaction.Height < height {
			return nil, 0, fmt.Errorf("unexpected transaction height order: %d < %d", transaction.Height, height)
		}

		height = transaction.Height
	}

	return g.Transactions, height, nil
}
