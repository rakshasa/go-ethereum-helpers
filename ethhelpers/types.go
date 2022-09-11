package ethhelpers

import (
	"math/big"
)

func BigIntAsInt64(v *big.Int) (int64, bool) {
	if v == nil || !v.IsInt64() {
		return 0, false
	}

	return v.Int64(), true
}

func BigIntAsInt64OrZero(v *big.Int) int64 {
	id, _ := BigIntAsInt64(v)
	return id
}

func BigIntAsUint64(v *big.Int) (uint64, bool) {
	if v == nil || !v.IsUint64() {
		return 0, false
	}

	return v.Uint64(), true
}

func BigIntAsUint64OrZero(v *big.Int) uint64 {
	id, _ := BigIntAsUint64(v)
	return id
}
