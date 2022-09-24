package ethtesting

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
)

var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr    = crypto.PubkeyToAddress(testKey.PublicKey)
	testBalance = big.NewInt(2e15)
)

var genesis = &core.Genesis{
	Config: params.AllEthashProtocolChanges,

	Alloc: core.GenesisAlloc{
		testAddr: {
			Balance: testBalance,
		},
	},
	ExtraData: []byte("test genesis"),
	Timestamp: 9000,
	BaseFee:   big.NewInt(params.InitialBaseFee),
}

func newTestBackend(t *testing.T) (*node.Node, []*types.Block) {
	blocks, _, err := GenerateTestChain(genesis, &SimpleBlockGenerator{
		ExpectedHeight: 10,
		ExtraGenerator: func(int, *core.BlockGen) []byte { return []byte("test backend") },
		Transactions: []TransactionWithHeight{
			{1, types.MustSignNewTx(testKey, types.LatestSigner(genesis.Config), &types.LegacyTx{
				Nonce:    0,
				Value:    big.NewInt(12),
				GasPrice: big.NewInt(params.InitialBaseFee),
				Gas:      params.TxGas,
				To:       &common.Address{2},
			})},
			{2, types.MustSignNewTx(testKey, types.LatestSigner(genesis.Config), &types.LegacyTx{
				Nonce:    1,
				Value:    big.NewInt(8),
				GasPrice: big.NewInt(params.InitialBaseFee),
				Gas:      params.TxGas,
				To:       &common.Address{2},
			})},
		},
	})
	if !assert.NoError(t, err) {
		t.Fatalf("could not generate test chain")
	}

	backend, err := NewTestBackend(
		node.Config{},
		genesis,
		blocks,
	)
	if !assert.NoError(t, err) {
		t.Fatalf("could not generate test backend")
	}

	return backend, blocks
}

func TestGenerateBackend(t *testing.T) {
	commit := PendingLogHandlerForTesting(t, log.Root())
	defer commit()

	backend, _ := newTestBackend(t)
	defer backend.Close()

	// TODO: Verify tx count and block numbers.

	rpcClient, _ := backend.Attach()
	defer rpcClient.Close()

	ec := ethclient.NewClient(rpcClient)

	id, err := ec.ChainID(context.Background())
	if assert.NoError(t, err) {
		return
	}

	assert.NotNil(t, id)
	assert.NotEqual(t, 0, id.Cmp(params.AllEthashProtocolChanges.ChainID))
}
