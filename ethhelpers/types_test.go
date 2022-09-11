package ethhelpers

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigInt(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		fn   func(string)
	}{
		{
			"nil",
			func(name string) {
				var v *big.Int

				assert.Equal(int64(0), BigIntAsInt64OrZero(v), name)
				assert.Equal(uint64(0), BigIntAsUint64OrZero(v), name)

				i64, ok := BigIntAsInt64(v)
				assert.Equal(int64(0), i64, name)
				assert.False(ok, name)

				u64, ok := BigIntAsUint64(v)
				assert.Equal(uint64(0), u64, name)
				assert.False(ok, name)
			},
		}, {
			"big.NewInt(0)",
			func(name string) {
				v := big.NewInt(0)

				assert.Equal(int64(0), BigIntAsInt64OrZero(v), name)
				assert.Equal(uint64(0), BigIntAsUint64OrZero(v), name)

				i64, ok := BigIntAsInt64(v)
				assert.Equal(int64(0), i64, name)
				assert.True(ok, name)

				u64, ok := BigIntAsUint64(v)
				assert.Equal(uint64(0), u64, name)
				assert.True(ok, name)
			},
		}, {
			"big.NewInt(1)",
			func(name string) {
				v := big.NewInt(1)

				assert.Equal(int64(1), BigIntAsInt64OrZero(v), name)
				assert.Equal(uint64(1), BigIntAsUint64OrZero(v), name)

				i64, ok := BigIntAsInt64(v)
				assert.Equal(int64(1), i64, name)
				assert.True(ok, name)

				u64, ok := BigIntAsUint64(v)
				assert.Equal(uint64(1), u64, name)
				assert.True(ok, name)
			},
		}, {
			"big.NewInt(1<<62)",
			func(name string) {
				v := new(big.Int).SetBit(big.NewInt(0), 62, 1)

				assert.Equal(int64(1<<62), BigIntAsInt64OrZero(v), name)
				assert.Equal(uint64(1<<62), BigIntAsUint64OrZero(v), name)

				i64, ok := BigIntAsInt64(v)
				assert.Equal(int64(1<<62), i64, name)
				assert.True(ok, name)

				u64, ok := BigIntAsUint64(v)
				assert.Equal(uint64(1<<62), u64, name)
				assert.True(ok, name)
			},
		}, {
			"big.NewInt(1<<63)",
			func(name string) {
				v := new(big.Int).SetBit(big.NewInt(0), 63, 1)

				assert.Equal(int64(0), BigIntAsInt64OrZero(v), name)
				assert.Equal(uint64(1<<63), BigIntAsUint64OrZero(v), name)

				i64, ok := BigIntAsInt64(v)
				assert.Equal(int64(0), i64, name)
				assert.False(ok, name)

				u64, ok := BigIntAsUint64(v)
				assert.Equal(uint64(1<<63), u64, name)
				assert.True(ok, name)
			},
		}, {
			"big.NewInt(1<<64)",
			func(name string) {
				v := new(big.Int).SetBit(big.NewInt(0), 64, 1)

				assert.Equal(int64(0), BigIntAsInt64OrZero(v), name)
				assert.Equal(uint64(0), BigIntAsUint64OrZero(v), name)

				i64, ok := BigIntAsInt64(v)
				assert.Equal(int64(0), i64, name)
				assert.False(ok, name)

				u64, ok := BigIntAsUint64(v)
				assert.Equal(uint64(0), u64, name)
				assert.False(ok, name)
			},
		},
	}

	for idx, test := range tests {
		test.fn(fmt.Sprintf("%d: %s", idx, test.name))
	}
}
