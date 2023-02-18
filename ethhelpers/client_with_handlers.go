package ethhelpers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ClientCaller struct {
	name string
	args []interface{}
	call func(context.Context, Client) error
}

func (c ClientCaller) Name() string {
	return c.name
}

func (c ClientCaller) Args() []interface{} {
	return c.args
}

func (c ClientCaller) Call(ctx context.Context, client Client) error {
	return c.call(ctx, client)
}

type clientWithHandlers struct {
	defaultHandler func(context.Context, ClientCaller) error
}

// NewClientWithHandlers creates a new client with custom handlers.
//
// The handlers cannot modify the content of the arguments or results, except
// for overriding the error returned.
//
// Handlers should return nil if it has not changed the error.
func NewClientWithDefaultHandler(defaultHandler func(context.Context, ClientCaller) error) Client {
	return &clientWithHandlers{
		defaultHandler: defaultHandler,
	}
}

func (c *clientWithHandlers) handle(ctx context.Context, callInfo ClientCaller) error {
	switch {
	case c.defaultHandler != nil:
		return c.defaultHandler(ctx, callInfo)
	default:
		return fmt.Errorf("no handler for %s", callInfo.Name())
	}
}

// BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
func (c *clientWithHandlers) BlockByHash(ctx context.Context, hash common.Hash) (r *types.Block, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "BlockByHash",
		args: []interface{}{hash},
		call: func(ctx context.Context, client Client) error {
			r, e = client.BlockByHash(ctx, hash)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
func (c *clientWithHandlers) BlockByNumber(ctx context.Context, number *big.Int) (r *types.Block, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "BlockByNumber",
		args: []interface{}{number},
		call: func(ctx context.Context, client Client) error {
			r, e = client.BlockByNumber(ctx, number)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
func (c *clientWithHandlers) HeaderByHash(ctx context.Context, hash common.Hash) (r *types.Header, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "HeaderByHash",
		args: []interface{}{hash},
		call: func(ctx context.Context, client Client) error {
			r, e = client.HeaderByHash(ctx, hash)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
func (c *clientWithHandlers) HeaderByNumber(ctx context.Context, number *big.Int) (r *types.Header, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "HeaderByNumber",
		args: []interface{}{number},
		call: func(ctx context.Context, client Client) error {
			r, e = client.HeaderByNumber(ctx, number)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
func (c *clientWithHandlers) TransactionCount(ctx context.Context, blockHash common.Hash) (r uint, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "TransactionCount",
		args: []interface{}{blockHash},
		call: func(ctx context.Context, client Client) error {
			r, e = client.TransactionCount(ctx, blockHash)
			return e
		},
	}); err != nil {
		return 0, err
	}

	return
}

// TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
func (c *clientWithHandlers) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (r *types.Transaction, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "TransactionInBlock",
		args: []interface{}{blockHash, index},
		call: func(ctx context.Context, client Client) error {
			r, e = client.TransactionInBlock(ctx, blockHash, index)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
func (c *clientWithHandlers) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (r ethereum.Subscription, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "SubscribeNewHead",
		args: []interface{}{ch},
		call: func(ctx context.Context, client Client) error {
			r, e = client.SubscribeNewHead(ctx, ch)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
func (c *clientWithHandlers) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (r *big.Int, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "BalanceAt",
		args: []interface{}{account, blockNumber},
		call: func(ctx context.Context, client Client) error {
			r, e = client.BalanceAt(ctx, account, blockNumber)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
func (c *clientWithHandlers) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) (r []byte, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "StorageAt",
		args: []interface{}{account, key, blockNumber},
		call: func(ctx context.Context, client Client) error {
			r, e = client.StorageAt(ctx, account, key, blockNumber)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
func (c *clientWithHandlers) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) (r []byte, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "CodeAt",
		args: []interface{}{account, blockNumber},
		call: func(ctx context.Context, client Client) error {
			r, e = client.CodeAt(ctx, account, blockNumber)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
func (c *clientWithHandlers) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (r uint64, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "NonceAt",
		args: []interface{}{account, blockNumber},
		call: func(ctx context.Context, client Client) error {
			r, e = client.NonceAt(ctx, account, blockNumber)
			return e
		},
	}); err != nil {
		return 0, err
	}

	return
}

// SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)
func (c *clientWithHandlers) SyncProgress(ctx context.Context) (r *ethereum.SyncProgress, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "SyncProgress",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			r, e = client.SyncProgress(ctx)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
func (c *clientWithHandlers) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) (r []byte, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "CallContract",
		args: []interface{}{call, blockNumber},
		call: func(ctx context.Context, client Client) error {
			r, e = client.CallContract(ctx, call, blockNumber)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
func (c *clientWithHandlers) EstimateGas(ctx context.Context, call ethereum.CallMsg) (r uint64, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "EstimateGas",
		args: []interface{}{call},
		call: func(ctx context.Context, client Client) error {
			r, e = client.EstimateGas(ctx, call)
			return e
		},
	}); err != nil {
		return 0, err
	}

	return
}

// SuggestGasPrice(ctx context.Context) (*big.Int, error)
func (c *clientWithHandlers) SuggestGasPrice(ctx context.Context) (r *big.Int, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "SuggestGasPrice",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			r, e = client.SuggestGasPrice(ctx)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
func (c *clientWithHandlers) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (r []types.Log, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "FilterLogs",
		args: []interface{}{q},
		call: func(ctx context.Context, client Client) error {
			r, e = client.FilterLogs(ctx, q)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
func (c *clientWithHandlers) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (r ethereum.Subscription, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "SubscribeFilterLogs",
		args: []interface{}{q, ch},
		call: func(ctx context.Context, client Client) error {
			r, e = client.SubscribeFilterLogs(ctx, q, ch)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error)
func (c *clientWithHandlers) PendingCallContract(ctx context.Context, call ethereum.CallMsg) (r []byte, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "PendingCallContract",
		args: []interface{}{call},
		call: func(ctx context.Context, client Client) error {
			r, e = client.PendingCallContract(ctx, call)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// SubscribePendingTransactions(ctx context.Context, ch chan<- *types.Transaction) (ethereum.Subscription, error)

// PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error)
func (c *clientWithHandlers) PendingBalanceAt(ctx context.Context, account common.Address) (r *big.Int, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "PendingBalanceAt",
		args: []interface{}{account},
		call: func(ctx context.Context, client Client) error {
			r, e = client.PendingBalanceAt(ctx, account)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error)
func (c *clientWithHandlers) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) (r []byte, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "PendingStorageAt",
		args: []interface{}{account, key},
		call: func(ctx context.Context, client Client) error {
			r, e = client.PendingStorageAt(ctx, account, key)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
func (c *clientWithHandlers) PendingCodeAt(ctx context.Context, account common.Address) (r []byte, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "PendingCodeAt",
		args: []interface{}{account},
		call: func(ctx context.Context, client Client) error {
			r, e = client.PendingCodeAt(ctx, account)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
func (c *clientWithHandlers) PendingNonceAt(ctx context.Context, account common.Address) (r uint64, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "PendingNonceAt",
		args: []interface{}{account},
		call: func(ctx context.Context, client Client) error {
			r, e = client.PendingNonceAt(ctx, account)
			return e
		},
	}); err != nil {
		return 0, err
	}

	return
}

// PendingTransactionCount(ctx context.Context) (uint, error)
func (c *clientWithHandlers) PendingTransactionCount(ctx context.Context) (r uint, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "PendingTransactionCount",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			r, e = client.PendingTransactionCount(ctx)
			return e
		},
	}); err != nil {
		return 0, err
	}

	return
}

// TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error)
func (c *clientWithHandlers) TransactionByHash(ctx context.Context, txHash common.Hash) (r *types.Transaction, isPending bool, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "TransactionByHash",
		args: []interface{}{txHash},
		call: func(ctx context.Context, client Client) error {
			r, isPending, e = client.TransactionByHash(ctx, txHash)
			return e
		},
	}); err != nil {
		return nil, false, err
	}

	return
}

// TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
func (c *clientWithHandlers) TransactionReceipt(ctx context.Context, txHash common.Hash) (r *types.Receipt, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "TransactionReceipt",
		args: []interface{}{txHash},
		call: func(ctx context.Context, client Client) error {
			r, e = client.TransactionReceipt(ctx, txHash)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// SendTransaction(ctx context.Context, tx *types.Transaction) error
func (c *clientWithHandlers) SendTransaction(ctx context.Context, tx *types.Transaction) (e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "SendTransaction",
		args: []interface{}{tx},
		call: func(ctx context.Context, client Client) error {
			e = client.SendTransaction(ctx, tx)
			return e
		},
	}); err != nil {
		return err
	}

	return
}

// BlockNumber(ctx context.Context) (uint64, error)
func (c *clientWithHandlers) BlockNumber(ctx context.Context) (r uint64, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "BlockNumber",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			r, e = client.BlockNumber(ctx)
			return e
		},
	}); err != nil {
		return 0, err
	}

	return
}

// CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error)
func (c *clientWithHandlers) CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) (r []byte, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "CallContractAtHash",
		args: []interface{}{msg, blockHash},
		call: func(ctx context.Context, client Client) error {
			r, e = client.CallContractAtHash(ctx, msg, blockHash)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// ChainID(ctx context.Context) (*big.Int, error)
func (c *clientWithHandlers) ChainID(ctx context.Context) (r *big.Int, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "ChainID",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			r, e = client.ChainID(ctx)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// Close()
func (c *clientWithHandlers) Close() {
	if err := c.handle(context.Background(), ClientCaller{
		name: "Close",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			client.Close()
			return nil
		},
	}); err != nil {
		return
	}
}

// FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
func (c *clientWithHandlers) FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (r *ethereum.FeeHistory, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "FeeHistory",
		args: []interface{}{blockCount, lastBlock, rewardPercentiles},
		call: func(ctx context.Context, client Client) error {
			r, e = client.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// NetworkID(ctx context.Context) (*big.Int, error)
func (c *clientWithHandlers) NetworkID(ctx context.Context) (r *big.Int, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "NetworkID",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			r, e = client.NetworkID(ctx)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}

// PeerCount(ctx context.Context) (uint64, error)
func (c *clientWithHandlers) PeerCount(ctx context.Context) (r uint64, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "PeerCount",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			r, e = client.PeerCount(ctx)
			return e
		},
	}); err != nil {
		return 0, err
	}

	return
}

// SuggestGasTipCap(ctx context.Context) (*big.Int, error)
func (c *clientWithHandlers) SuggestGasTipCap(ctx context.Context) (r *big.Int, e error) {
	if err := c.handle(ctx, ClientCaller{
		name: "SuggestGasTipCap",
		args: []interface{}{},
		call: func(ctx context.Context, client Client) error {
			r, e = client.SuggestGasTipCap(ctx)
			return e
		},
	}); err != nil {
		return nil, err
	}

	return
}
