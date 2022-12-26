package ethhelpers_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethhelpers"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientWithDefaultHandler(t *testing.T) {
	type testArgs struct {
		ctx    context.Context
		client ethhelpers.ClientWithHandlers
		mock   *mock.Mock
	}

	tests := []struct {
		name string
		call func(*testing.T, testArgs)
	}{
		{
			name: "BlockByHash returns success",
			call: func(t *testing.T, args testArgs) {
				hash := common.HexToHash("0x1234")
				expectedResult := &types.Block{}

				args.mock.On("BlockByHash", args.ctx, hash).Return(expectedResult, nil).Once()

				r, err := args.client.BlockByHash(args.ctx, hash)
				assert.NoError(t, err)
				assert.Same(t, expectedResult, r)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			client := ethtesting.NewClientWithMock()

			test.call(t, testArgs{
				ctx,
				ethhelpers.NewClientWithDefaultHandler(func(ci ethhelpers.ClientCallInfo) error {
					return ci.Call(ci.Ctx, client)
				}),
				client.Mock(),
			})

			client.Mock().AssertExpectations(t)
		})
	}
}
