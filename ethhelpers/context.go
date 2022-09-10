package ethhelpers

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type clientContextKey struct{}
type configContextKey struct{}
type rpcClientContextKey struct{}

// ContextWithClient creates a new context which contains an ethclient client.
//
// The context will return the client when calling ClientFromContext.
func ContextWithClient(ctx context.Context, client *ethclient.Client) context.Context {
	return context.WithValue(ctx, clientContextKey{}, client)
}

// ContextWithClients creates a new context which contains both the
// RPC client and a newly created ethclient client.
//
// The context will return the clients when calling ClientFromContext and
// RPCClientFromContext.
func ContextWithClients(ctx context.Context, rpcClient *rpc.Client) context.Context {
	ctx = context.WithValue(ctx, clientContextKey{}, ethclient.NewClient(rpcClient))
	ctx = context.WithValue(ctx, rpcClientContextKey{}, rpcClient)
	return ctx
}

// ContextWithConfig creates a new context which contains an ethhelpers config.
//
// The context will return the config when calling ConfigFromContext.
func ContextWithConfig(ctx context.Context, config Config) context.Context {
	return context.WithValue(ctx, configContextKey{}, config)
}

// ContextWithRPCClient creates a new context which contains an RPC client.
//
// The context will return the client when calling RPCClientFromContext.
func ContextWithRPCClient(ctx context.Context, rpcClient *rpc.Client) context.Context {
	return context.WithValue(ctx, rpcClientContextKey{}, rpcClient)
}

// ClientFromContext retrieves an ethclient client from the context, if any.
func ClientFromContext(ctx context.Context) (*ethclient.Client, bool) {
	c, ok := ctx.Value(clientContextKey{}).(*ethclient.Client)
	return c, ok
}

// ConfigFromContext retrieves an ethhelpers config from the context, if any.
func ConfigFromContext(ctx context.Context) (Config, bool) {
	c, ok := ctx.Value(configContextKey{}).(Config)
	return c, ok
}

// RPCClientFromContext retrieves an RPC client from the context, if any.
func RPCClientFromContext(ctx context.Context) (*rpc.Client, bool) {
	c, ok := ctx.Value(rpcClientContextKey{}).(*rpc.Client)
	return c, ok
}
