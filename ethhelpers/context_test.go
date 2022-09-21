package ethhelpers

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

func TestClientFromContext(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		fn   func(string)
	}{
		{
			"client : empty context",
			func(name string) {
				c, ok := ClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"limited client : empty context",
			func(name string) {
				c, ok := LimitedClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"rpc client : empty context",
			func(name string) {
				c, ok := RPCClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"with clients from rpc client",
			func(name string) {
				stored := &rpc.Client{}
				ctx := ContextWithClientsFromRPCClient(context.Background(), stored)

				e, ok := ClientFromContext(ctx)
				assert.NotNil(e, name)
				assert.True(ok, name)
				c, ok := ConfigFromContext(ctx)
				assert.Equal(Config{}, c, name)
				assert.False(ok, name)
				l, ok := LimitedClientFromContext(ctx)
				assert.NotNil(l, name)
				assert.True(ok, name)
				r, ok := RPCClientFromContext(ctx)
				assert.Same(stored, r, name)
				assert.True(ok, name)
			},
		}, {
			"with ethclient client",
			func(name string) {
				stored := &ethclient.Client{}
				ctx := ContextWithClient(context.Background(), stored)

				e, ok := ClientFromContext(ctx)
				assert.Same(stored, e, name)
				assert.True(ok, name)
				c, ok := ConfigFromContext(ctx)
				assert.Equal(Config{}, c, name)
				assert.False(ok, name)
				l, ok := LimitedClientFromContext(ctx)
				assert.NotNil(l, name)
				assert.True(ok, name)
				r, ok := RPCClientFromContext(ctx)
				assert.Nil(r, name)
				assert.False(ok, name)
			},
		}, {
			"with rpc client",
			func(name string) {
				stored := &rpc.Client{}
				ctx := ContextWithRPCClient(context.Background(), stored)

				e, ok := ClientFromContext(ctx)
				assert.Nil(e, name)
				assert.False(ok, name)
				c, ok := ConfigFromContext(ctx)
				assert.Equal(Config{}, c, name)
				assert.False(ok, name)
				l, ok := LimitedClientFromContext(ctx)
				assert.Nil(l, name)
				assert.False(ok, name)
				r, ok := RPCClientFromContext(ctx)
				assert.Same(stored, r, name)
				assert.True(ok, name)
			},
		}, {
			"with simulated backend client",
			func(name string) {
				stored := &backends.SimulatedBackend{}
				ctx := ContextWithLimitedClient(context.Background(), stored)

				e, ok := ClientFromContext(ctx)
				assert.Nil(e, name)
				assert.False(ok, name)
				c, ok := ConfigFromContext(ctx)
				assert.Equal(Config{}, c, name)
				assert.False(ok, name)
				l, ok := LimitedClientFromContext(ctx)
				assert.NotNil(l, name)
				assert.True(ok, name)
				r, ok := RPCClientFromContext(ctx)
				assert.Nil(r, name)
				assert.False(ok, name)
			},
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}

func TestConfigFromContext(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		fn   func(string)
	}{
		{
			"empty context",
			func(name string) {
				c, ok := RPCClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"with config",
			func(name string) {
				stored := Config{Endpoint: "test"}
				ctx := ContextWithConfig(context.Background(), stored)

				e, ok := ClientFromContext(ctx)
				assert.Nil(e, name)
				assert.False(ok, name)
				c, ok := ConfigFromContext(ctx)
				assert.Equal(stored, c, name)
				assert.True(ok, name)
				r, ok := RPCClientFromContext(ctx)
				assert.Nil(r, name)
				assert.False(ok, name)
			},
		}, {
			"with config and clients",
			func(name string) {
				storedConfig := Config{
					Endpoint: "test",
					ChainId:  big.NewInt(5),
				}
				storedRpcClient := &rpc.Client{}

				ctx := context.Background()
				ctx = ContextWithConfig(ctx, storedConfig)
				ctx = ContextWithClientsFromRPCClient(ctx, storedRpcClient)

				e, ok := ClientFromContext(ctx)
				assert.NotNil(e, name)
				assert.True(ok, name)
				c, ok := ConfigFromContext(ctx)
				assert.Equal(storedConfig, c, name)
				assert.True(ok, name)
				r, ok := RPCClientFromContext(ctx)
				assert.Same(storedRpcClient, r, name)
				assert.True(ok, name)
			},
		}, {
			"with contract",
			func(name string) {
				// TODO: Add tests with empty Contracts.
				config := Config{Contracts: NewContractContainer()}
				ctx := ContextWithConfig(context.Background(), config)

				c, ok := ContractFromConfigInContext(ctx, 1)
				assert.Nil(c, name)
				assert.False(ok, name)

				c = ContractOrNilFromConfigInContext(ctx, 1)
				assert.Nil(c, name)

				stored := &testContract{}
				config.Contracts.Put(1, stored)

				c, ok = ContractFromConfigInContext(ctx, 1)
				assert.Equal(stored, c, name)
				assert.True(ok, name)

				c = ContractOrNilFromConfigInContext(ctx, 1)
				assert.Equal(stored, c, name)
			},
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}
