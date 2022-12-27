package ethhelpers_test

import (
	"context"
	"fmt"
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
		client ethhelpers.Client
		mock   *mock.Mock
	}

	preCallError := fmt.Errorf("pre call error")
	postCallError := fmt.Errorf("post call error")
	callError := fmt.Errorf("call error")

	tests := []struct {
		name          string
		call          func(*testing.T, testArgs)
		preCallError  error
		postCallError error
	}{
		// BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
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
		}, {
			name: "BlockByHash call returns failure",
			call: func(t *testing.T, args testArgs) {
				hash := common.HexToHash("0x1234")

				args.mock.On("BlockByHash", args.ctx, hash).Return(nil, callError).Once()

				r, err := args.client.BlockByHash(args.ctx, hash)
				assert.Same(t, callError, err)
				assert.Nil(t, r)
			},
		}, {
			name: "BlockByHash pre handler returns failure",
			call: func(t *testing.T, args testArgs) {
				hash := common.HexToHash("0x1234")

				r, err := args.client.BlockByHash(args.ctx, hash)
				assert.Same(t, callError, err)
				assert.Nil(t, r)
			},
			preCallError: preCallError,
		}, {
			name: "BlockByHash post handler returns failure",
			call: func(t *testing.T, args testArgs) {
				hash := common.HexToHash("0x1234")
				expectedResult := &types.Block{}

				args.mock.On("BlockByHash", args.ctx, hash).Return(expectedResult, nil).Once()

				r, err := args.client.BlockByHash(args.ctx, hash)
				assert.Same(t, postCallError, err)
				assert.Nil(t, r)
			},
			postCallError: postCallError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			client := ethtesting.NewClientWithMock()

			test.call(t, testArgs{
				ctx,
				ethhelpers.NewClientWithDefaultHandler(func(ctx context.Context, caller ethhelpers.ClientCaller) (e error) {
					if test.preCallError != nil {
						return test.preCallError
					}

					e = caller.Call(ctx, client)

					if test.postCallError != nil {
						return test.postCallError
					}

					return
				}),
				client.Mock(),
			})

			client.Mock().AssertExpectations(t)
		})
	}
}
