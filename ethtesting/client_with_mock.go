package ethtesting

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/stretchr/testify/mock"
)

type CanceledMockCallOption struct{}

func CanceledMockCall() CanceledMockCallOption {
	return CanceledMockCallOption{}
}

type PassthroughMockCallOption struct{}

func PassthroughMockCall() PassthroughMockCallOption {
	return PassthroughMockCallOption{}
}

type ClientWithMock interface {
	ethhelpers.Client
	Mock() *mock.Mock
	Test(t mock.TestingT)
}

type clientWithMock struct {
	mu     sync.Mutex
	client ethhelpers.Client
	mock   mock.Mock
	test   mock.TestingT
}

// Newethhelpers.Client creates a new client with *testing.Mock.
func NewClientWithMock() ClientWithMock {
	return &clientWithMock{}
}

// NewClientWithMockAndClient creates a new client with *testing.Mock and an
// underlying client.
//
// Using PassthroughMockCall() as the assigned return value for mocked calls will pass
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

// Test sets the test struct variable of the mock object.
func (c *clientWithMock) Test(t mock.TestingT) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.mock.Test(t)
	c.test = t
}

func (c *clientWithMock) fail(format string, args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.test == nil {
		panic(fmt.Sprintf(format, args...))
	}

	c.test.Errorf(format, args...)
	c.test.FailNow()
}

// BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
func (c *clientWithMock) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	values := c.mock.MethodCalled("BlockByHash", ctx, hash)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.BlockByHash(ctx, hash)
	case *types.Block:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
func (c *clientWithMock) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	values := c.mock.MethodCalled("BlockByNumber", ctx, number)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.BlockByNumber(ctx, number)
	case *types.Block:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
func (c *clientWithMock) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	values := c.mock.MethodCalled("HeaderByHash", ctx, hash)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.HeaderByHash(ctx, hash)
	case *types.Header:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
func (c *clientWithMock) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	values := c.mock.MethodCalled("HeaderByNumber", ctx, number)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.HeaderByNumber(ctx, number)
	case *types.Header:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
func (c *clientWithMock) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	values := c.mock.MethodCalled("TransactionCount", ctx, blockHash)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return 0, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return 0, nil
		}
		return c.client.TransactionCount(ctx, blockHash)
	case uint:
		return v0, values.Error(1)
	case nil:
		return 0, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return 0, nil
	}
}

// TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
func (c *clientWithMock) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	values := c.mock.MethodCalled("TransactionInBlock", ctx, blockHash, index)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.TransactionInBlock(ctx, blockHash, index)
	case *types.Transaction:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
func (c *clientWithMock) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	values := c.mock.MethodCalled("SubscribeNewHead", ctx, ch)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.SubscribeNewHead(ctx, ch)
	case ethereum.Subscription:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
func (c *clientWithMock) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	values := c.mock.MethodCalled("BalanceAt", ctx, account, blockNumber)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.BalanceAt(ctx, account, blockNumber)
	case *big.Int:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
func (c *clientWithMock) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	values := c.mock.MethodCalled("StorageAt", ctx, account, key, blockNumber)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.StorageAt(ctx, account, key, blockNumber)
	case []byte:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
func (c *clientWithMock) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	values := c.mock.MethodCalled("CodeAt", ctx, account, blockNumber)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.CodeAt(ctx, account, blockNumber)
	case []byte:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
func (c *clientWithMock) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	values := c.mock.MethodCalled("NonceAt", ctx, account, blockNumber)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return 0, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return 0, nil
		}
		return c.client.NonceAt(ctx, account, blockNumber)
	case uint64:
		return v0, values.Error(1)
	case nil:
		return 0, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return 0, nil
	}
}

// SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
func (c *clientWithMock) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	values := c.mock.MethodCalled("SyncProgress", ctx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.SyncProgress(ctx)
	case *ethereum.SyncProgress:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
func (c *clientWithMock) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	values := c.mock.MethodCalled("CallContract", ctx, call, blockNumber)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.CallContract(ctx, call, blockNumber)
	case []byte:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
func (c *clientWithMock) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	values := c.mock.MethodCalled("EstimateGas", ctx, call)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return 0, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return 0, nil
		}
		return c.client.EstimateGas(ctx, call)
	case uint64:
		return v0, values.Error(1)
	case nil:
		return 0, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return 0, nil
	}
}

// SuggestGasPrice(ctx context.Context) (*big.Int, error)
func (c *clientWithMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	values := c.mock.MethodCalled("SuggestGasPrice", ctx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.SuggestGasPrice(ctx)
	case *big.Int:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
func (c *clientWithMock) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	values := c.mock.MethodCalled("FilterLogs", ctx, q)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.FilterLogs(ctx, q)
	case []types.Log:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
func (c *clientWithMock) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	values := c.mock.MethodCalled("SubscribeFilterLogs", ctx, q, ch)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.SubscribeFilterLogs(ctx, q, ch)
	case ethereum.Subscription:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error)
func (c *clientWithMock) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	values := c.mock.MethodCalled("PendingCallContract", ctx, call)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.PendingCallContract(ctx, call)
	case []byte:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
func (c *clientWithMock) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	values := c.mock.MethodCalled("PendingBalanceAt", ctx, account)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.PendingBalanceAt(ctx, account)
	case *big.Int:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error)
func (c *clientWithMock) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	values := c.mock.MethodCalled("PendingStorageAt", ctx, account, key)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.PendingStorageAt(ctx, account, key)
	case []byte:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
func (c *clientWithMock) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	values := c.mock.MethodCalled("PendingCodeAt", ctx, account)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.PendingCodeAt(ctx, account)
	case []byte:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
func (c *clientWithMock) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	values := c.mock.MethodCalled("PendingNonceAt", ctx, account)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return 0, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return 0, nil
		}
		return c.client.PendingNonceAt(ctx, account)
	case uint64:
		return v0, values.Error(1)
	case nil:
		return 0, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return 0, nil
	}
}

// PendingTransactionCount(ctx context.Context) (uint, error)
func (c *clientWithMock) PendingTransactionCount(ctx context.Context) (uint, error) {
	values := c.mock.MethodCalled("PendingTransactionCount", ctx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return 0, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return 0, nil
		}
		return c.client.PendingTransactionCount(ctx)
	case uint:
		return v0, values.Error(1)
	case nil:
		return 0, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return 0, nil
	}
}

// TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error)
func (c *clientWithMock) TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	values := c.mock.MethodCalled("TransactionByHash", ctx, txHash)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, false, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, false, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, false, nil
		}
		return c.client.TransactionByHash(ctx, txHash)
	case *types.Transaction:
		if v1, ok := values.Get(1).(bool); ok {
			return v0, v1, values.Error(2)
		} else {
			c.fail("unexpected mock return type: %T", v1)
			return nil, false, nil
		}
	case nil:
		if v1, ok := values.Get(1).(bool); ok {
			return nil, v1, values.Error(2)
		} else {
			c.fail("unexpected mock return type: %T", v1)
			return nil, false, nil
		}
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, false, nil
	}
}

// TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
func (c *clientWithMock) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	values := c.mock.MethodCalled("TransactionReceipt", ctx, txHash)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.TransactionReceipt(ctx, txHash)
	case *types.Receipt:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// SendTransaction(ctx context.Context, tx *types.Transaction) error
func (c *clientWithMock) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	values := c.mock.MethodCalled("SendTransaction", ctx, tx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil
		}
		return c.client.SendTransaction(ctx, tx)
	case error:
		return v0
	case nil:
		return nil
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil
	}
}

// BlockNumber(ctx context.Context) (uint64, error)
func (c *clientWithMock) BlockNumber(ctx context.Context) (uint64, error) {
	values := c.mock.MethodCalled("BlockNumber", ctx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return 0, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return 0, nil
		}
		return c.client.BlockNumber(ctx)
	case uint64:
		return v0, values.Error(1)
	case nil:
		return 0, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return 0, nil
	}
}

// CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error)
func (c *clientWithMock) CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	values := c.mock.MethodCalled("CallContractAtHash", ctx, msg, blockHash)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.CallContractAtHash(ctx, msg, blockHash)
	case []byte:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// ChainID(ctx context.Context) (*big.Int, error)
func (c *clientWithMock) ChainID(ctx context.Context) (*big.Int, error) {
	values := c.mock.MethodCalled("ChainID", ctx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.ChainID(ctx)
	case *big.Int:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// Close()
func (c *clientWithMock) Close() {
	values := c.mock.MethodCalled("Close")

	switch v0 := values.Get(0).(type) {
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return
		}
		c.client.Close()
		return
	case nil:
		return
	default:
		c.fail("unexpected mock return type: %T", v0)
		return
	}
}

// FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
func (c *clientWithMock) FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error) {
	values := c.mock.MethodCalled("FeeHistory", ctx, blockCount, lastBlock, rewardPercentiles)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
	case *ethereum.FeeHistory:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// NetworkID(ctx context.Context) (*big.Int, error)
func (c *clientWithMock) NetworkID(ctx context.Context) (*big.Int, error) {
	values := c.mock.MethodCalled("NetworkID", ctx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.NetworkID(ctx)
	case *big.Int:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}

// PeerCount(ctx context.Context) (uint64, error)
func (c *clientWithMock) PeerCount(ctx context.Context) (uint64, error) {
	values := c.mock.MethodCalled("PeerCount", ctx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return 0, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return 0, nil
		}
		return c.client.PeerCount(ctx)
	case uint64:
		return v0, values.Error(1)
	case nil:
		return 0, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return 0, nil
	}
}

// SuggestGasTipCap(ctx context.Context) (*big.Int, error)
func (c *clientWithMock) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	values := c.mock.MethodCalled("SuggestGasTipCap", ctx)

	switch v0 := values.Get(0).(type) {
	case CanceledMockCallOption:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.fail("mock call was canceled but context was not canceled")
			return nil, nil
		}
	case PassthroughMockCallOption:
		if c.client == nil {
			c.fail("client is nil")
			return nil, nil
		}
		return c.client.SuggestGasTipCap(ctx)
	case *big.Int:
		return v0, values.Error(1)
	case nil:
		return nil, values.Error(1)
	default:
		c.fail("unexpected mock return type: %T", v0)
		return nil, nil
	}
}
