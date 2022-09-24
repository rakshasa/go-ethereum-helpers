package ethtesting

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	MockPrivateKey1 = MustHexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	MockPrivateKey2 = MustHexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	MockPrivateKey3 = MustHexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	MockPrivateKey4 = MustHexToECDSA("e238eb8e04fee6511ab04c6dd3c89ce097b11f25d584863ac2b6d5b35b1847e4")
)

func MustGenerateKey() *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(fmt.Sprintf("could not generate key: %v", err))
	}

	return key
}

// MustHexToECDSA returns the result of HexToECDSA on hexkey, or panics if error.
func MustHexToECDSA(hexkey string) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		panic(fmt.Sprintf("%s : %v", hexkey, err))
	}

	return key
}
