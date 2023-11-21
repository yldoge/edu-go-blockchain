package core

import (
	"fmt"

	"github.com/yldoge/edu-go-blockchain/crypto"
)

type Transaction struct {
	Data []byte

	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}

func (tx *Transaction) Sign(pvk crypto.PrivateKey) error {
	sig, err := pvk.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.PublicKey = pvk.PublicKey()
	tx.Signature = sig
	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.PublicKey, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}
	return nil
}
