package ethtesting_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rakshasa/go-ethereum-helpers/ethtesting"
	"github.com/stretchr/testify/assert"
)

type mockSubscription struct{}

func (m *mockSubscription) Unsubscribe()      {}
func (m *mockSubscription) Err() <-chan error { return nil }

func TestClientWithMock(t *testing.T) {
	type testArgs struct {
		ctx    context.Context
		sim    *ethtesting.SimulatedBackendWithAccounts
		client ethtesting.ClientWithMock
	}

	tests := []struct {
		name string
		call func(*testing.T, testArgs)
	}{
		{
			name: "BlockNumber",
			call: func(t *testing.T, args testArgs) {
				args.client.Mock().On("BlockNumber", args.ctx).Return(uint64(123), fmt.Errorf("mocked error")).Once()
				args.client.Mock().On("BlockNumber", args.ctx).Return(ethtesting.WithoutMock()).Once()

				blockNumber, err := args.client.BlockNumber(args.ctx)
				assert.Error(t, err)
				assert.Equal(t, uint64(123), blockNumber)

				blockNumber, err = args.client.BlockNumber(args.ctx)
				assert.NoError(t, err)
				assert.Equal(t, uint64(1), blockNumber)
			},
		}, {
			name: "FilterLogs",
			call: func(t *testing.T, args testArgs) {
				q := ethereum.FilterQuery{
					FromBlock: big.NewInt(1),
				}
				fl := []types.Log{types.Log{}}

				args.client.Mock().On("FilterLogs", args.ctx, q).Return(fl, fmt.Errorf("mocked error")).Once()
				args.client.Mock().On("FilterLogs", args.ctx, q).Return(ethtesting.WithoutMock()).Once()

				logs, err := args.client.FilterLogs(args.ctx, q)
				assert.Error(t, err)
				assert.Equal(t, fl, logs)

				logs, err = args.client.FilterLogs(args.ctx, q)
				assert.NoError(t, err)
				assert.Equal(t, []types.Log{}, logs)
			},
		}, {
			name: "SubscribeFilterLogs",
			call: func(t *testing.T, args testArgs) {
				q := ethereum.FilterQuery{
					FromBlock: big.NewInt(1),
				}
				ch := make(chan types.Log)
				expectedSub := &mockSubscription{}

				args.client.Mock().On("SubscribeFilterLogs", args.ctx, q, (chan<- types.Log)(ch)).Return(expectedSub, fmt.Errorf("mocked error")).Once()
				args.client.Mock().On("SubscribeFilterLogs", args.ctx, q, (chan<- types.Log)(ch)).Return(ethtesting.WithoutMock()).Once()

				sub, err := args.client.SubscribeFilterLogs(args.ctx, q, ch)
				assert.Error(t, err)
				assert.Equal(t, expectedSub, sub)

				sub, err = args.client.SubscribeFilterLogs(args.ctx, q, ch)
				assert.NoError(t, err)
				assert.NotNil(t, sub)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			sim := ethtesting.NewSimulatedBackendWithAccounts(
				ethtesting.GenesisAccountWithPrivateKey{
					PrivateKey: ethtesting.MockPrivateKey1,
					GenesisAccount: core.GenesisAccount{
						Balance: big.NewInt(10_000_000_000_000_000),
					},
				},
				ethtesting.GenesisAccountWithPrivateKey{
					PrivateKey: ethtesting.MockPrivateKey2,
				},
			)
			defer sim.Backend.Close()

			sim.Backend.Commit()

			ctx := context.Background()
			client := ethtesting.NewClientWithMockAndClient(ethtesting.NewSimulatedClient(sim.Backend))

			test.call(t, testArgs{ctx, sim, client})

			client.Mock().AssertExpectations(t)
		})
	}
}
