package ethhelpers

import (
	"context"
	"fmt"
	"testing"

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
				assert.Nil(c, "empty context")
				assert.False(ok, "empty context")
			},
		}, {
			"ethclient : with client",
			func(name string) {
				stored := &rpc.Client{}
				c, ok := ClientFromContext(ContextWithClients(context.Background(), stored))
				assert.NotNil(c, name)
				assert.True(ok, name)
			},
		}, {
			"rpc : with client",
			func(name string) {
				stored := &rpc.Client{}
				c, ok := RPCClientFromContext(ContextWithClients(context.Background(), stored))
				assert.Same(stored, c, name)
				assert.True(ok, name)
			},
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}
