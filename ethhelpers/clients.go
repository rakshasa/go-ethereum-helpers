package ethhelpers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Client is an interface that holds the same methods as ethclient.Client.
type Client interface {
	//
	// Methods from go-ethereum interfaces:
	//

	// ethereum.ChainReader
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)

	// ethereum.ChainStateReader
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)

	// ethereum.ChainSyncReader
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)

	// ethereum.ContractCaller
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)

	// ethereum.GasEstimator
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)

	// ethereum.GasPricer
	SuggestGasPrice(ctx context.Context) (*big.Int, error)

	// ethereum.LogFilterer
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)

	// ethereum.PendingContractCaller
	PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error)

	// ethereum.PendingStateEventer (geth only)
	// SubscribePendingTransactions(ctx context.Context, ch chan<- *types.Transaction) (ethereum.Subscription, error)

	// ethereum.PendingStateReader
	PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
	PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	PendingTransactionCount(ctx context.Context) (uint, error)

	// ethereum.TransactionReader
	TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

	// ethereum.TransactionSender
	SendTransaction(ctx context.Context, tx *types.Transaction) error

	//
	// Methods from ethclient interfaces:
	//

	// BlockNumberReader
	BlockNumber(ctx context.Context) (uint64, error)

	//
	// Methods not included in any interface:
	//

	CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error)
	ChainID(ctx context.Context) (*big.Int, error)
	// TODO: Remove Close?
	Close()
	FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	PeerCount(ctx context.Context) (uint64, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

// TODO: Add client without subcription methods.

type BlockNumberReader interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type ChainReaderAndTransactionSender interface {
	ethereum.TransactionSender
	ethereum.ChainReader
}

type FilterLogsReader interface {
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}
