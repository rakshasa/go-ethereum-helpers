package ethtesting

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func MustHexToECDSA(hexkey string) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		panic(fmt.Sprintf("%s : %v", hexkey, err))
	}

	return key
}
