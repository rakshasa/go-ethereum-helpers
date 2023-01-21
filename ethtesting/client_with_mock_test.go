package ethtesting_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

type mockSubscription struct{}

func (m *mockSubscription) Unsubscribe()      {}
func (m *mockSubscription) Err() <-chan error { return nil }

func TestClientWithMock(t *testing.T) {
	type testInfo struct {
		name   string
		ctx    context.Context
		sim    *ethtesting.SimulatedBackendWithAccounts
		client ethtesting.ClientWithMock
	}

	type testCall1 struct {
		fn                func(ctx context.Context) (interface{}, error)
		args              []interface{}
		mockResult        interface{}
		defaultResult     interface{}
		defaultIsNil      bool
		passthroughResult interface{}
		passthroughNotNil bool
		passthroughError  bool
	}

	call1 := func(t *testing.T, info testInfo, call testCall1) {
		info.client.Mock().On(info.name, append([]interface{}{info.ctx}, call.args...)...).Return(call.mockResult, nil).Once()
		r, err := call.fn(info.ctx)
		assert.NoError(t, err)
		assert.Equal(t, call.mockResult, r)

		info.client.Mock().On(info.name, append([]interface{}{info.ctx}, call.args...)...).Return(call.defaultResult, fmt.Errorf("mocked error")).Once()
		r, err = call.fn(info.ctx)
		assert.Error(t, err)
		if call.defaultIsNil {
			assert.Nil(t, r)
		} else {
			assert.Equal(t, call.defaultResult, r)
		}

		info.client.Mock().On(info.name, append([]interface{}{info.ctx}, call.args...)...).Return(nil, fmt.Errorf("mocked error")).Once()
		r, err = call.fn(info.ctx)
		assert.Error(t, err)
		if call.defaultIsNil {
			assert.Nil(t, r)
		} else {
			assert.Equal(t, call.defaultResult, r)
		}

		info.client.Mock().On(info.name, append([]interface{}{info.ctx}, call.args...)...).Return(ethtesting.PassthroughMockCall()).Once()
		r, err = call.fn(info.ctx)
		if call.passthroughError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		if call.passthroughNotNil {
			assert.NotNil(t, r)
		} else {
			assert.Equal(t, call.passthroughResult, r)
		}

		canceledCtx, cancel := context.WithCancel(info.ctx)
		info.client.Mock().On(info.name, append([]interface{}{canceledCtx}, call.args...)...).Return(ethtesting.CanceledMockCall()).Once()
		cancel()
		r, err = call.fn(canceledCtx)
		assert.Error(t, err)
		if call.defaultIsNil {
			assert.Nil(t, r)
		} else {
			assert.Equal(t, call.defaultResult, r)
		}
	}

	tests := []struct {
		name string
		call func(*testing.T, testInfo)
	}{
		{
			name: "BlockByHash",
			call: func(t *testing.T, info testInfo) {
				call1(t, info, testCall1{
					fn: func(ctx context.Context) (interface{}, error) {
						return info.client.BlockByHash(ctx, common.HexToHash("0x1234"))
					},
					args:              []interface{}{common.HexToHash("0x1234")},
					mockResult:        &types.Block{},
					defaultIsNil:      true,
					passthroughResult: (*types.Block)(nil),
					passthroughError:  true,
				})
			},
		}, {
			name: "BlockNumber",
			call: func(t *testing.T, info testInfo) {
				call1(t, info, testCall1{
					fn: func(ctx context.Context) (interface{}, error) {
						return info.client.BlockNumber(ctx)
					},
					args:              []interface{}{},
					mockResult:        uint64(123),
					defaultResult:     uint64(0),
					passthroughResult: uint64(1),
					passthroughError:  false,
				})
			},
		}, {
			name: "FilterLogs",
			call: func(t *testing.T, info testInfo) {
				call1(t, info, testCall1{
					fn: func(ctx context.Context) (interface{}, error) {
						return info.client.FilterLogs(ctx, ethereum.FilterQuery{FromBlock: big.NewInt(1)})
					},
					args:              []interface{}{ethereum.FilterQuery{FromBlock: big.NewInt(1)}},
					mockResult:        []types.Log{types.Log{}},
					defaultIsNil:      true,
					passthroughResult: []types.Log{},
					passthroughError:  false,
				})
			},
		}, {
			name: "SubscribeFilterLogs",
			call: func(t *testing.T, info testInfo) {
				ch := make(chan types.Log)

				call1(t, info, testCall1{
					fn: func(ctx context.Context) (interface{}, error) {
						return info.client.SubscribeFilterLogs(ctx, ethereum.FilterQuery{FromBlock: big.NewInt(1)}, ch)
					},
					args:              []interface{}{ethereum.FilterQuery{FromBlock: big.NewInt(1)}, (chan<- types.Log)(ch)},
					mockResult:        &mockSubscription{},
					defaultIsNil:      true,
					passthroughError:  false,
					passthroughNotNil: true,
				})
			},
		},
	}

	// SendTransaction

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			sim, close := newDefaultSimulatedBackend(t)
			defer close()

			sim.Backend.Commit()

			ctx := context.Background()
			client := ethtesting.NewClientWithMockAndClient(ethtesting.NewSimulatedClient(sim.Backend))
			client.Test(t)

			test.call(t, testInfo{test.name, ctx, sim, client})

			client.Mock().AssertExpectations(t)
		})
	}
}
