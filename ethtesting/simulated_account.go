package ethtesting

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type SimulatedAccount struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
	Nonce      uint64
}

func NewSimulatedAccount(privateKey *ecdsa.PrivateKey) *SimulatedAccount {
	return NewSimulatedAccountWithNonce(privateKey, 0)
}

func NewSimulatedAccountWithNonce(privateKey *ecdsa.PrivateKey, nonce uint64) *SimulatedAccount {
	return &SimulatedAccount{
		PrivateKey: privateKey,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		Nonce:      nonce,
	}
}

// NextNonce increments SimulatedAccount.Nonce and returns its old value.
func (a *SimulatedAccount) NonceAndIncrement() (n uint64) {
	n = a.Nonce
	a.Nonce++
	return
}

func (a *SimulatedAccount) SendNewTransaction(ctx context.Context, client ChainReaderAndTransactionSender, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, data []byte) (*types.Transaction, error) {
	head, _ := client.HeaderByNumber(ctx, nil)
	gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(1))

	signedTx, err := types.SignTx(
		types.NewTransaction(nonce, to, amount, params.TxGas, gasPrice, data),
		types.HomesteadSigner{},
		a.PrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("could not sign transaction: %v", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("could not add transaction to pending block: %v", err)
	}

	return signedTx, nil
}
