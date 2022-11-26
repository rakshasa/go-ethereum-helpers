package ethtesting

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/stretchr/testify/mock"
)

type WithoutMockOption struct{}

func WithoutMock() WithoutMockOption {
	return WithoutMockOption{}
}

type ClientWithMock interface {
	ethereum.LogFilterer

	// ethereum.ChainReader
	// ethereum.ChainStateReader
	// ethereum.ChainSyncReader
	// ethereum.ContractCaller
	// ethereum.GasEstimator
	// ethereum.GasPricer
	// ethereum.PendingContractCaller
	// ethereum.PendingStateReader
	// ethereum.TransactionReader
	// ethereum.TransactionSender

	ethhelpers.BlockNumberReader

	// CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error)
	// ChainID(ctx context.Context) (*big.Int, error)
	// Close()
	// FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
	// NetworkID(ctx context.Context) (*big.Int, error)
	// PeerCount(ctx context.Context) (uint64, error)
	// SendTransaction(ctx context.Context, tx *types.Transaction) error
	// SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	// SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	// TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	// TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
}

type clientWithMock struct {
	client ethhelpers.Client
	mock   mock.Mock
}

func NewClientWithMock(client ethhelpers.Client) (ClientWithMock, *mock.Mock) {
	c := &clientWithMock{
		client: client,
	}
	return c, &c.mock
}

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
