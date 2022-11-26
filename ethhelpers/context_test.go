package ethhelpers_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

func TestContext_ClientFromContext(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		fn   func(string)
	}{
		{
			"client : empty context",
			func(name string) {
				c, ok := ethhelpers.ClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"rpc client : empty context",
			func(name string) {
				c, ok := ethhelpers.RPCClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"with clients from rpc client",
			func(name string) {
				stored := &rpc.Client{}
				ctx := ethhelpers.ContextWithClientsFromRPCClient(context.Background(), stored)

				e, ok := ethhelpers.ClientFromContext(ctx)
				assert.NotNil(e, name)
				assert.True(ok, name)
				c, ok := ethhelpers.ConfigFromContext(ctx)
				assert.Equal(ethhelpers.Config{}, c, name)
				assert.False(ok, name)
				r, ok := ethhelpers.RPCClientFromContext(ctx)
				assert.Same(stored, r, name)
				assert.True(ok, name)
			},
		}, {
			"with ethclient client",
			func(name string) {
				stored := &ethclient.Client{}
				ctx := ethhelpers.ContextWithClient(context.Background(), stored)

				e, ok := ethhelpers.ClientFromContext(ctx)
				assert.Same(stored, e, name)
				assert.True(ok, name)
				c, ok := ethhelpers.ConfigFromContext(ctx)
				assert.Equal(ethhelpers.Config{}, c, name)
				assert.False(ok, name)
				r, ok := ethhelpers.RPCClientFromContext(ctx)
				assert.Nil(r, name)
				assert.False(ok, name)
			},
		}, {
			"with rpc client",
			func(name string) {
				stored := &rpc.Client{}
				ctx := ethhelpers.ContextWithRPCClient(context.Background(), stored)

				e, ok := ethhelpers.ClientFromContext(ctx)
				assert.Nil(e, name)
				assert.False(ok, name)
				c, ok := ethhelpers.ConfigFromContext(ctx)
				assert.Equal(ethhelpers.Config{}, c, name)
				assert.False(ok, name)
				r, ok := ethhelpers.RPCClientFromContext(ctx)
				assert.Same(stored, r, name)
				assert.True(ok, name)
			},
		}, {
			"with simulated backend using ethtesting.SimulatedClient",
			func(name string) {
				stored := ethtesting.NewSimulatedClient(&backends.SimulatedBackend{})
				ctx := ethhelpers.ContextWithClient(context.Background(), stored)

				e, ok := ethhelpers.ClientFromContext(ctx)
				assert.Equal(stored, e, name)
				assert.True(ok, name)
				c, ok := ethhelpers.ConfigFromContext(ctx)
				assert.Equal(ethhelpers.Config{}, c, name)
				assert.False(ok, name)
				r, ok := ethhelpers.RPCClientFromContext(ctx)
				assert.Nil(r, name)
				assert.False(ok, name)
			},
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}

func TestContext_ConfigFromContext(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		fn   func(string)
	}{
		{
			"empty context",
			func(name string) {
				c, ok := ethhelpers.RPCClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"with config",
			func(name string) {
				stored := ethhelpers.Config{Endpoint: "test"}
				ctx := ethhelpers.ContextWithConfig(context.Background(), stored)

				e, ok := ethhelpers.ClientFromContext(ctx)
				assert.Nil(e, name)
				assert.False(ok, name)
				c, ok := ethhelpers.ConfigFromContext(ctx)
				assert.Equal(stored, c, name)
				assert.True(ok, name)
				r, ok := ethhelpers.RPCClientFromContext(ctx)
				assert.Nil(r, name)
				assert.False(ok, name)
			},
		}, {
			"with config and clients",
			func(name string) {
				storedConfig := ethhelpers.Config{
					Endpoint: "test",
					ChainId:  big.NewInt(5),
				}
				storedRpcClient := &rpc.Client{}

				ctx := context.Background()
				ctx = ethhelpers.ContextWithConfig(ctx, storedConfig)
				ctx = ethhelpers.ContextWithClientsFromRPCClient(ctx, storedRpcClient)

				e, ok := ethhelpers.ClientFromContext(ctx)
				assert.NotNil(e, name)
				assert.True(ok, name)
				c, ok := ethhelpers.ConfigFromContext(ctx)
				assert.Equal(storedConfig, c, name)
				assert.True(ok, name)
				r, ok := ethhelpers.RPCClientFromContext(ctx)
				assert.Same(storedRpcClient, r, name)
				assert.True(ok, name)
			},
		}, {
			"with contract",
			func(name string) {
				// TODO: Add tests with empty Contracts.
				config := ethhelpers.Config{Contracts: ethhelpers.NewContractContainer()}
				ctx := ethhelpers.ContextWithConfig(context.Background(), config)

				c, ok := ethhelpers.ContractFromConfigInContext(ctx, 1)
				assert.Nil(c, name)
				assert.False(ok, name)

				c = ethhelpers.ContractOrNilFromConfigInContext(ctx, 1)
				assert.Nil(c, name)

				stored := &testContract{}
				config.Contracts.Put(1, stored)

				c, ok = ethhelpers.ContractFromConfigInContext(ctx, 1)
				assert.Equal(stored, c, name)
				assert.True(ok, name)

				c = ethhelpers.ContractOrNilFromConfigInContext(ctx, 1)
				assert.Equal(stored, c, name)
			},
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}
