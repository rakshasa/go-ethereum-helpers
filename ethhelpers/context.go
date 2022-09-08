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

// ContextWithClient creates a new context which contains an ethclient client.
//
// The context will return the client when calling ClientFromContext.
func ContextWithClient(ctx context.Context, client *ethclient.Client) context.Context {
	c, ok := ctx.Value(clientContextKey{}).(clientContext)
	if !ok {
		c = clientContext{}
	}

	c.client = client

	return context.WithValue(ctx, clientContextKey{}, c)
}

// ContextWithClients creates a new context which contains both the
// RPC client and a newly created ethclient client.
//
// The context will return the clients when calling ClientFromContext and
// RPCClientFromContext.
func ContextWithClients(ctx context.Context, rpcClient *rpc.Client) context.Context {
	c, ok := ctx.Value(clientContextKey{}).(clientContext)
	if !ok {
		c = clientContext{}
	}

	c.client = ethclient.NewClient(rpcClient)
	c.rpcClient = rpcClient

	return context.WithValue(ctx, clientContextKey{}, c)
}

// ContextWithRPCClient creates a new context which contains an RPC client.
//
// The context will return the client when calling RPCClientFromContext.
func ContextWithRPCClient(ctx context.Context, rpcClient *rpc.Client) context.Context {
	c, ok := ctx.Value(clientContextKey{}).(clientContext)
	if !ok {
		c = clientContext{}
	}

	c.rpcClient = rpcClient

	return context.WithValue(ctx, clientContextKey{}, c)
}

// RPCClientFromContext retrieves an ethclient client from the context, if any.
func ClientFromContext(ctx context.Context) (*ethclient.Client, bool) {
	c, ok := ctx.Value(clientContextKey{}).(clientContext)
	if !ok || c.client == nil {
		return nil, false
	}

	return c.client, true
}

// RPCClientFromContext retrieves an RPC client from the context, if any.
func RPCClientFromContext(ctx context.Context) (*rpc.Client, bool) {
	c, ok := ctx.Value(clientContextKey{}).(clientContext)
	if !ok || c.rpcClient == nil {
		return nil, false
	}

	return c.rpcClient, true
}
