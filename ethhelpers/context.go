package ethhelpers

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type clientContextKey struct{}

type clientContext struct {
	client    *ethclient.Client
	rpcClient *rpc.Client
}

// ContextWithClients stores in a new context both the RPC client and
// the ethclient client created from it.
//
// The context will return clients for both ClientFromContext and RPCClientFromContext.
func ContextWithClients(ctx context.Context, rpcClient *rpc.Client) context.Context {
	c, ok := ctx.Value(clientContextKey{}).(clientContext)
	if !ok {
		c = clientContext{}
	}

	c.client = ethclient.NewClient(rpcClient)
	c.rpcClient = rpcClient

	return context.WithValue(ctx, clientContextKey{}, c)
}

// RPCClientFromContext retrieves an ethclient client from the context, if any.
func ClientFromContext(ctx context.Context) (*ethclient.Client, bool) {
	c, ok := ctx.Value(clientContextKey{}).(clientContext)
	if !ok {
		return nil, false
	}

	return c.client, true
}

// RPCClientFromContext retrieves an RPC client from the context, if any.
func RPCClientFromContext(ctx context.Context) (*rpc.Client, bool) {
	c, ok := ctx.Value(clientContextKey{}).(clientContext)
	if !ok {
		return nil, false
	}

	return c.rpcClient, true
}
