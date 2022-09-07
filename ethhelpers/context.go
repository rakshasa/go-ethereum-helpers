package ethhelpers

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type rpcClientContextKey struct{}

// ContextWithRPCClient stores the RPC client in a new context.
func ContextWithRPCClient(ctx context.Context, rpcClient *rpc.Client) context.Context {
	return context.WithValue(ctx, rpcClientContextKey{}, rpcClient)
}

// ClientFromContext retrieves a new client created from the RPC client from the context, if any.
func ClientFromContext(ctx context.Context) (*ethclient.Client, bool) {
	c, ok := ctx.Value(rpcClientContextKey{}).(*rpc.Client)
	if !ok {
		return nil, false
	}

	return ethclient.NewClient(c), true
}

// RPCClientFromContext retrieves an RPC client from the context, if any.
func RPCClientFromContext(ctx context.Context) (c *rpc.Client, ok bool) {
	c, ok = ctx.Value(rpcClientContextKey{}).(*rpc.Client)
	return
}
