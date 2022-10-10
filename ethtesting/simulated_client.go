package ethtesting

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SimulatedClient interface {
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
	FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
	NetworkID(ctx context.Context) (*big.Int, error)
	PeerCount(ctx context.Context) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error)
}

type simulatedClient struct {
	*backends.SimulatedBackend
}

func NewSimulatedClient(sim *backends.SimulatedBackend) SimulatedClient {
	return &simulatedClient{
		SimulatedBackend: sim,
	}
}

func (c *simulatedClient) BlockNumber(ctx context.Context) (uint64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		v := c.SimulatedBackend.Blockchain().CurrentHeader().Number
		if v == nil || !v.IsUint64() {
			return 0, fmt.Errorf("not a valid uint64 number")
		}

		return v.Uint64(), nil
	}
}

func (c *simulatedClient) CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *simulatedClient) ChainID(ctx context.Context) (*big.Int, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return SimulatedChainID(), nil
	}
}

func (c *simulatedClient) Close() {
	_ = c.SimulatedBackend.Close()
}

func (c *simulatedClient) FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *simulatedClient) NetworkID(ctx context.Context) (*big.Int, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *simulatedClient) PeerCount(ctx context.Context) (uint64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (c *simulatedClient) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *simulatedClient) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *simulatedClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	return 0, fmt.Errorf("not implemented")
}

func (c *simulatedClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	return nil, fmt.Errorf("not implemented")
}
