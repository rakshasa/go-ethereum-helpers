package ethhelpers

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type configContextKey struct{}

// ContextWithConfig creates a new context which contains an ethhelpers config.
//
// The context will return the config when calling ConfigFromContext.
func ContextWithConfig(ctx context.Context, config Config) context.Context {
	return context.WithValue(ctx, configContextKey{}, config)
}

// ConfigFromContext retrieves an ethhelpers.Config from the context, if any.
func ConfigFromContext(ctx context.Context) (Config, bool) {
	c, ok := ctx.Value(configContextKey{}).(Config)
	return c, ok
}

// ContractFromConfigInContext retrieves a Contract from the
// ContractContainer in the Config stored in the context, if any.
func ContractFromConfigInContext(ctx context.Context, key interface{}) (Contract, bool) {
	config, ok := ctx.Value(configContextKey{}).(Config)
	if !ok {
		return nil, false
	}

	c, ok := config.Contracts.Get(key)
	return c, ok
}

// Same as ContractFromConfigInContext, except it returns a nil object
// if not present.
//
//	contractHelper, ok := ethhelpers.ContractOrNilFromConfigInContext(ctx, MyContractKey{}).(*MyContract)
//	if !ok {
//	  return fmt.Errorf("missing my contract in context")
//	}
//
//	contract, err := contractHelper.ContractFromContext(ctx)
//	if err != nil {
//	  return fmt.Errorf("failed to create my contract: %v", err)
//	}
func ContractOrNilFromConfigInContext(ctx context.Context, key interface{}) Contract {
	c, _ := ContractFromConfigInContext(ctx, key)
	return c
}

type clientContextKey struct{}
type rpcClientContextKey struct{}

// ContextWithClient creates a new context which contains an ethhelpers.Client.
//
// The context will return the client when calling ClientFromContext
// and other compatible methods.
//
// Note that there can only be one client type in addition to the RPC
// client stored in the context, other client interface variants are
// stored using the same context key.
func ContextWithClient(ctx context.Context, client Client) context.Context {
	return context.WithValue(ctx, clientContextKey{}, client)
}

// ContextWithClients creates a new context which contains both the
// RPC client and a newly created ethclient client.
//
// The context will return the clients when calling ClientFromContext and
// RPCClientFromContext.
func ContextWithClientsFromRPCClient(ctx context.Context, rpcClient *rpc.Client) context.Context {
	ctx = context.WithValue(ctx, clientContextKey{}, ethclient.NewClient(rpcClient))
	ctx = context.WithValue(ctx, rpcClientContextKey{}, rpcClient)
	return ctx
}

// ContextWithRPCClient creates a new context which contains an RPC client.
//
// The context will return the client when calling RPCClientFromContext.
func ContextWithRPCClient(ctx context.Context, rpcClient *rpc.Client) context.Context {
	return context.WithValue(ctx, rpcClientContextKey{}, rpcClient)
}

// ClientFromContext retrieves an interface implementing
// ethhelpers.Client from the context, if any.
func ClientFromContext(ctx context.Context) (Client, bool) {
	c, ok := ctx.Value(clientContextKey{}).(Client)
	return c, ok
}

// RPCClientFromContext retrieves an *rpc.Client from the context, if any.
func RPCClientFromContext(ctx context.Context) (*rpc.Client, bool) {
	c, ok := ctx.Value(rpcClientContextKey{}).(*rpc.Client)
	return c, ok
}
