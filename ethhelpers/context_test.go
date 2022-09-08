package ethhelpers

import (
	"context"
	"fmt"
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
			"rpc : empty context",
			func(name string) {
				c, ok := RPCClientFromContext(context.Background())
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"ethclient : with clients",
			func(name string) {
				stored := &rpc.Client{}
				c, ok := ClientFromContext(ContextWithClients(context.Background(), stored))
				assert.NotNil(c, name)
				assert.True(ok, name)
			},
		}, {
			"rpc : with clients",
			func(name string) {
				stored := &rpc.Client{}
				c, ok := RPCClientFromContext(ContextWithClients(context.Background(), stored))
				assert.Same(stored, c, name)
				assert.True(ok, name)
			},
		}, {
			"ethclient : with ethclient client",
			func(name string) {
				stored := &ethclient.Client{}
				c, ok := ClientFromContext(ContextWithClient(context.Background(), stored))
				assert.Same(stored, c, name)
				assert.True(ok, name)
			},
		}, {
			"rpc : with ethclient client",
			func(name string) {
				stored := &ethclient.Client{}
				c, ok := RPCClientFromContext(ContextWithClient(context.Background(), stored))
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"ethclient : with rpc client",
			func(name string) {
				stored := &rpc.Client{}
				c, ok := ClientFromContext(ContextWithRPCClient(context.Background(), stored))
				assert.Nil(c, name)
				assert.False(ok, name)
			},
		}, {
			"rpc : with rpc client",
			func(name string) {
				stored := &rpc.Client{}
				c, ok := RPCClientFromContext(ContextWithRPCClient(context.Background(), stored))
				assert.Same(stored, c, name)
				assert.True(ok, name)
			},
		}, {
			"ethclient and rpc : with clients",
			func(name string) {
				storedRpc := &rpc.Client{}
				storedEthclient := &ethclient.Client{}

				ctx := context.Background()
				ctx = ContextWithClient(ctx, storedEthclient)
				ctx = ContextWithRPCClient(ctx, storedRpc)

				ethClient, ok := ClientFromContext(ctx)
				assert.Same(storedEthclient, ethClient, name)
				assert.True(ok, name)

				rpcClient, ok := RPCClientFromContext(ctx)
				assert.Same(storedRpc, rpcClient, name)
				assert.True(ok, name)
			},
		}, {
			"ethclient and rpc : with clients and overwrite",
			func(name string) {
				storedRpc := &rpc.Client{}
				storedEthclient := &ethclient.Client{}

				ctx := context.Background()
				ctx = ContextWithClient(ctx, storedEthclient)
				ctx = ContextWithRPCClient(ctx, storedRpc)
				ctx = ContextWithClient(ctx, nil)

				ethClient, ok := ClientFromContext(ctx)
				assert.Nil(ethClient, name)
				assert.False(ok, name)

				rpcClient, ok := RPCClientFromContext(ctx)
				assert.Same(storedRpc, rpcClient, name)
				assert.True(ok, name)

				ctx = ContextWithRPCClient(ctx, nil)

				ethClient, ok = ClientFromContext(ctx)
				assert.Nil(ethClient, name)
				assert.False(ok, name)

				rpcClient, ok = RPCClientFromContext(ctx)
				assert.Nil(rpcClient, name)
				assert.False(ok, name)
			},
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}
