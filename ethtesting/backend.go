package ethtesting

import (
	"fmt"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
)

func GenerateTestChain(genesis *core.Genesis, blockGenerator BlockGenerator) ([]*types.Block, []types.Receipts, error) {
	generator, height, err := blockGenerator.Generator()
	if err != nil {
		return nil, nil, fmt.Errorf("could not create block generator: %v", err)
	}

	// From 1.10.25:

	db := rawdb.NewMemoryDatabase()

	genesisBlock, err := genesis.Commit(db)
	if err != nil {
		return nil, nil, fmt.Errorf("could not commit genesis block: %v", err)
	}

	blocks, receipts := core.GenerateChain(
		genesis.Config,
		genesisBlock,
		ethash.NewFaker(),
		db,
		height,
		generator,
	)

	return append([]*types.Block{genesisBlock}, blocks...), receipts, nil

	// From master:

	// _, blocks, _ := core.GenerateChainWithGenesis(genesis, ethash.NewFaker(), 2, generate)

	// return append([]*types.Block{genesis.ToBlock()}, blocks...)
}

func NewTestBackend(genesis *core.Genesis, blocks []*types.Block) (*node.Node, error) {
	backend, err := node.New(&node.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't create new node: %v", err)
	}

	config := &ethconfig.Config{
		Genesis: genesis,
		Ethash: ethash.Config{
			PowMode: ethash.ModeFake,
		},
	}

	service, err := eth.New(backend, config)
	if err != nil {
		return nil, fmt.Errorf("can't create new ethereum service: %v", err)
	}

	if err := backend.Start(); err != nil {
		return nil, fmt.Errorf("can't start test node: %v", err)
	}
	if _, err := service.BlockChain().InsertChain(blocks[1:]); err != nil {
		return nil, fmt.Errorf("can't import test blocks: %v", err)
	}

	return backend, nil
}
