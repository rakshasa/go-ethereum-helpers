package ethhelpers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TODO: Consider a better name for this as it is based on SimulatedBackend.

type LimitedClient interface {
	ethereum.ChainReader
	ethereum.ChainStateReader
	// ethereum.ChainSyncReader
	ethereum.ContractCaller
	ethereum.GasEstimator
	ethereum.GasPricer
	ethereum.LogFilterer
	ethereum.PendingContractCaller
	// ethereum.PendingStateReader
	ethereum.TransactionReader
	ethereum.TransactionSender

	// Methods not included in any of the go-ethereum interfaces:

	// BlockNumber(ctx context.Context) (uint64, error)
	// CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error)
	// ChainID(ctx context.Context) (*big.Int, error)
	// Close()
	// FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
	// NetworkID(ctx context.Context) (*big.Int, error)
	// PeerCount(ctx context.Context) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
}

type Client interface {
	ethereum.ChainReader
	ethereum.ChainStateReader
	ethereum.ChainSyncReader
	ethereum.ContractCaller
	ethereum.GasEstimator
	ethereum.GasPricer
	ethereum.LogFilterer
	ethereum.PendingContractCaller
	ethereum.PendingStateReader
	ethereum.TransactionReader
	ethereum.TransactionSender

	// Methods not included in any of the go-ethereum interfaces:

	BlockNumber(ctx context.Context) (uint64, error)
	CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error)
	ChainID(ctx context.Context) (*big.Int, error)
	Close()
	// FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	PeerCount(ctx context.Context) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
}
