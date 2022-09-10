package ethhelpers

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

func TestFromContext(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		fn   func(string)
	}{
		{
			"ethclient : empty context",
			func(name string) {
				c, ok := ClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"config : empty context",
			func(name string) {
				c, ok := RPCClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"rpc : empty context",
			func(name string) {
				c, ok := RPCClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"with clients",
			func(name string) {
				stored := &rpc.Client{}
				ctx := ContextWithClients(context.Background(), stored)

				e, ok := ClientFromContext(ctx)
				assert.NotNil(e, name)
				assert.True(ok, name)
				c, ok := ConfigFromContext(ctx)
				assert.Equal(Config{}, c, name)
				assert.False(ok, name)
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
				r, ok := RPCClientFromContext(ctx)
				assert.Nil(r, name)
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
				r, ok := RPCClientFromContext(ctx)
				assert.Same(stored, r, name)
				assert.True(ok, name)
			},
		}, {
			"with config",
			func(name string) {
				storedConfig := Config{
					Endpoint: "test",
					ChainId:  big.NewInt(5),
				}
				storedRpcClient := &rpc.Client{}

				ctx := context.Background()
				ctx = ContextWithConfig(ctx, storedConfig)
				ctx = ContextWithClients(ctx, storedRpcClient)

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
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}
