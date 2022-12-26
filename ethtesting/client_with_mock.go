package ethtesting

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/stretchr/testify/mock"
)

type WithoutMockOption struct{}

func WithoutMock() WithoutMockOption {
	return WithoutMockOption{}
}

type ClientWithMock interface {
	ethhelpers.Client
	Mock() *mock.Mock
}

type clientWithMock struct {
	client ethhelpers.Client
	mock   mock.Mock
}

// Newethhelpers.Client creates a new client with *testing.Mock.
func NewClientWithMock() ClientWithMock {
	return &clientWithMock{}
}

// NewClientWithMockAndClient creates a new client with *testing.Mock and an
// underlying client.
//
// Using WithoutMock() as the assigned return value for mocked calls will pass
// the call to the underlying client.
func NewClientWithMockAndClient(client ethhelpers.Client) ClientWithMock {
	return &clientWithMock{
		client: client,
	}
}

// Mock returns the testify mock object.
func (c *clientWithMock) Mock() *mock.Mock {
	return &c.mock
}

// BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
func (c *clientWithMock) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	values := c.mock.MethodCalled("BlockByHash", ctx, hash)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.BlockByHash(ctx, hash)
	}

	return values.Get(0).(*types.Block), values.Error(1)
}

// BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
func (c *clientWithMock) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	values := c.mock.MethodCalled("BlockByNumber", ctx, number)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.BlockByNumber(ctx, number)
	}

	return values.Get(0).(*types.Block), values.Error(1)
}

// HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
func (c *clientWithMock) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	values := c.mock.MethodCalled("HeaderByHash", ctx, hash)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.HeaderByHash(ctx, hash)
	}

	return values.Get(0).(*types.Header), values.Error(1)
}

// HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
func (c *clientWithMock) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	values := c.mock.MethodCalled("HeaderByNumber", ctx, number)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.HeaderByNumber(ctx, number)
	}

	return values.Get(0).(*types.Header), values.Error(1)
}

// TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
func (c *clientWithMock) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	values := c.mock.MethodCalled("TransactionCount", ctx, blockHash)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.TransactionCount(ctx, blockHash)
	}

	return values.Get(0).(uint), values.Error(1)
}

// TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
func (c *clientWithMock) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	values := c.mock.MethodCalled("TransactionInBlock", ctx, blockHash, index)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.TransactionInBlock(ctx, blockHash, index)
	}

	return values.Get(0).(*types.Transaction), values.Error(1)
}

// SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
func (c *clientWithMock) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	values := c.mock.MethodCalled("SubscribeNewHead", ctx, ch)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.SubscribeNewHead(ctx, ch)
	}

	return values.Get(0).(ethereum.Subscription), values.Error(1)
}

// BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
func (c *clientWithMock) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	values := c.mock.MethodCalled("BalanceAt", ctx, account, blockNumber)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.BalanceAt(ctx, account, blockNumber)
	}

	return values.Get(0).(*big.Int), values.Error(1)
}

// StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
func (c *clientWithMock) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	values := c.mock.MethodCalled("StorageAt", ctx, account, key, blockNumber)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.StorageAt(ctx, account, key, blockNumber)
	}

	return values.Get(0).([]byte), values.Error(1)
}

// CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
func (c *clientWithMock) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	values := c.mock.MethodCalled("CodeAt", ctx, account, blockNumber)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.CodeAt(ctx, account, blockNumber)
	}

	return values.Get(0).([]byte), values.Error(1)
}

// NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
func (c *clientWithMock) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	values := c.mock.MethodCalled("NonceAt", ctx, account, blockNumber)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.NonceAt(ctx, account, blockNumber)
	}

	return values.Get(0).(uint64), values.Error(1)
}

// SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
func (c *clientWithMock) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	values := c.mock.MethodCalled("SyncProgress", ctx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.SyncProgress(ctx)
	}

	return values.Get(0).(*ethereum.SyncProgress), values.Error(1)
}

// CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
func (c *clientWithMock) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	values := c.mock.MethodCalled("CallContract", ctx, call, blockNumber)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.CallContract(ctx, call, blockNumber)
	}

	return values.Get(0).([]byte), values.Error(1)
}

// EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
func (c *clientWithMock) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	values := c.mock.MethodCalled("EstimateGas", ctx, call)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.EstimateGas(ctx, call)
	}

	return values.Get(0).(uint64), values.Error(1)
}

// SuggestGasPrice(ctx context.Context) (*big.Int, error)
func (c *clientWithMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	values := c.mock.MethodCalled("SuggestGasPrice", ctx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.SuggestGasPrice(ctx)
	}

	return values.Get(0).(*big.Int), values.Error(1)
}

// FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
func (c *clientWithMock) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	values := c.mock.MethodCalled("FilterLogs", ctx, q)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.FilterLogs(ctx, q)
	}

	return values.Get(0).([]types.Log), values.Error(1)
}

// SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
func (c *clientWithMock) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	values := c.mock.MethodCalled("SubscribeFilterLogs", ctx, q, ch)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.SubscribeFilterLogs(ctx, q, ch)
	}

	return values.Get(0).(ethereum.Subscription), values.Error(1)
}

// PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error)
func (c *clientWithMock) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	values := c.mock.MethodCalled("PendingCallContract", ctx, call)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.PendingCallContract(ctx, call)
	}

	return values.Get(0).([]byte), values.Error(1)
}

// PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
func (c *clientWithMock) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	values := c.mock.MethodCalled("PendingBalanceAt", ctx, account)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.PendingBalanceAt(ctx, account)
	}

	return values.Get(0).(*big.Int), values.Error(1)
}

// PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error)
func (c *clientWithMock) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	values := c.mock.MethodCalled("PendingStorageAt", ctx, account, key)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.PendingStorageAt(ctx, account, key)
	}

	return values.Get(0).([]byte), values.Error(1)
}

// PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
func (c *clientWithMock) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	values := c.mock.MethodCalled("PendingCodeAt", ctx, account)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.PendingCodeAt(ctx, account)
	}

	return values.Get(0).([]byte), values.Error(1)
}

// PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
func (c *clientWithMock) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	values := c.mock.MethodCalled("PendingNonceAt", ctx, account)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.PendingNonceAt(ctx, account)
	}

	return values.Get(0).(uint64), values.Error(1)
}

// PendingTransactionCount(ctx context.Context) (uint, error)
func (c *clientWithMock) PendingTransactionCount(ctx context.Context) (uint, error) {
	values := c.mock.MethodCalled("PendingTransactionCount", ctx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.PendingTransactionCount(ctx)
	}

	return values.Get(0).(uint), values.Error(1)
}

// TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error)
func (c *clientWithMock) TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	values := c.mock.MethodCalled("TransactionByHash", ctx, txHash)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.TransactionByHash(ctx, txHash)
	}

	return values.Get(0).(*types.Transaction), values.Get(1).(bool), values.Error(2)
}

// TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
func (c *clientWithMock) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	values := c.mock.MethodCalled("TransactionReceipt", ctx, txHash)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.TransactionReceipt(ctx, txHash)
	}

	return values.Get(0).(*types.Receipt), values.Error(1)
}

// SendTransaction(ctx context.Context, tx *types.Transaction) error
func (c *clientWithMock) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	values := c.mock.MethodCalled("SendTransaction", ctx, tx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.SendTransaction(ctx, tx)
	}

	return values.Error(0)
}

// BlockNumber(ctx context.Context) (uint64, error)
func (c *clientWithMock) BlockNumber(ctx context.Context) (uint64, error) {
	values := c.mock.MethodCalled("BlockNumber", ctx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.BlockNumber(ctx)
	}

	return values.Get(0).(uint64), values.Error(1)
}

// CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error)
func (c *clientWithMock) CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	values := c.mock.MethodCalled("CallContractAtHash", ctx, msg, blockHash)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.CallContractAtHash(ctx, msg, blockHash)
	}

	return values.Get(0).([]byte), values.Error(1)
}

// ChainID(ctx context.Context) (*big.Int, error)
func (c *clientWithMock) ChainID(ctx context.Context) (*big.Int, error) {
	values := c.mock.MethodCalled("ChainID", ctx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.ChainID(ctx)
	}

	return values.Get(0).(*big.Int), values.Error(1)
}

// Close()
func (c *clientWithMock) Close() {
	c.mock.MethodCalled("Close")
}

// FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
func (c *clientWithMock) FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error) {
	values := c.mock.MethodCalled("FeeHistory", ctx, blockCount, lastBlock, rewardPercentiles)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
	}

	return values.Get(0).(*ethereum.FeeHistory), values.Error(1)
}

// NetworkID(ctx context.Context) (*big.Int, error)
func (c *clientWithMock) NetworkID(ctx context.Context) (*big.Int, error) {
	values := c.mock.MethodCalled("NetworkID", ctx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.NetworkID(ctx)
	}

	return values.Get(0).(*big.Int), values.Error(1)
}

// PeerCount(ctx context.Context) (uint64, error)
func (c *clientWithMock) PeerCount(ctx context.Context) (uint64, error) {
	values := c.mock.MethodCalled("PeerCount", ctx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.PeerCount(ctx)
	}

	return values.Get(0).(uint64), values.Error(1)
}

// SuggestGasTipCap(ctx context.Context) (*big.Int, error)
func (c *clientWithMock) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	values := c.mock.MethodCalled("SuggestGasTipCap", ctx)

	if _, ok := values.Get(0).(WithoutMockOption); ok {
		if c.client == nil {
			panic("client is nil")
		}
		return c.client.SuggestGasTipCap(ctx)
	}

	return values.Get(0).(*big.Int), values.Error(1)
}
