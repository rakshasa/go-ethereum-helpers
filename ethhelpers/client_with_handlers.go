package ethhelpers

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ClientCallInfo struct {
	Ctx  context.Context
	Name string
	Call func(context.Context, Client) error
	Args []interface{}
}

type clientWithHandlers struct {
	defaultHandler func(ClientCallInfo) error
}

// TODO: Temporary.
type ClientWithHandlers interface {
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
}

// NewClientWithHandlers creates a new client with custom handlers.
//
// The handlers cannot modify the content of the arguments or results, except by
// overriding the error.
func NewClientWithDefaultHandler(defaultHandler func(ClientCallInfo) error) ClientWithHandlers {
	return &clientWithHandlers{
		defaultHandler: defaultHandler,
	}
}

func (c *clientWithHandlers) handle(callInfo ClientCallInfo) error {
	switch {
	case c.defaultHandler != nil:
		return c.defaultHandler(callInfo)
	default:
		return fmt.Errorf("no handler for %s", callInfo.Name)
	}
}

func (c *clientWithHandlers) BlockByHash(ctx context.Context, hash common.Hash) (r *types.Block, e error) {
	if err := c.handle(ClientCallInfo{
		Ctx:  ctx,
		Name: "BlockByHash",
		Call: func(ctx context.Context, client Client) error {
			r, e = client.BlockByHash(ctx, hash)
			return e
		},
		Args: []interface{}{hash},
	}); err != nil {
		return nil, err
	}

	return
}
