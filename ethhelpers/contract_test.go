package ethhelpers

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type testContract struct {
}

func (c *testContract) ChainId() *big.Int {
	return big.NewInt(1)
}

func (c *testContract) Address() common.Address {
	return common.HexToAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
}

func TestContractContainer(t *testing.T) {
	assert := assert.New(t)

	type key1 struct{}
	type key2 struct{}

	tests := []struct {
		name string
		fn   func(string)
	}{
		{
			"empty",
			func(name string) {
				c := NewContractContainer()

				contract, ok := c.Get(key1{})
				assert.Nil(contract, name)
				assert.False(ok, name)
			},
		}, {
			"put",
			func(name string) {
				c := NewContractContainer()
				expectedContract1 := &testContract{}
				expectedContract2 := &testContract{}

				contract, ok := c.Get(key1{})
				assert.Nil(contract, name)
				assert.False(ok, name)

				c.Put(key1{}, expectedContract1)

				contract, ok = c.Get(key1{})
				assert.Equal(expectedContract1, contract, name)
				assert.True(ok, name)

				contract, ok = c.Get(key2{})
				assert.Nil(contract, name)
				assert.False(ok, name)

				c.Put(key2{}, expectedContract2)

				contract, ok = c.Get(key1{})
				assert.Equal(expectedContract1, contract, name)
				assert.True(ok, name)

				contract, ok = c.Get(key2{})
				assert.Equal(expectedContract2, contract, name)
				assert.True(ok, name)
			},
		}, {
			"delete",
			func(name string) {
				c := NewContractContainer()
				expectedContract1 := &testContract{}
				expectedContract2 := &testContract{}

				c.Put(key1{}, expectedContract1)
				c.Put(key2{}, expectedContract2)

				c.Delete(key1{})

				contract, ok := c.Get(key1{})
				assert.Nil(contract, name)
				assert.False(ok, name)

				contract, ok = c.Get(key2{})
				assert.Equal(expectedContract2, contract, name)
				assert.True(ok, name)
			},
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}
