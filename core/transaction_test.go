package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yldoge/edu-go-blockchain/crypto"
)

func TestTransactionSign(t *testing.T) {
	pvk := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(pvk))
	assert.NotNil(t, tx.Signature)
}

func TestTransactionVerify(t *testing.T) {
	pvk := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(pvk))
	assert.Nil(t, tx.Verify())

	// invalid public key
	tx.PublicKey = crypto.GeneratePrivateKey().PublicKey()
	assert.NotNil(t, tx.Verify())

	// change data
	tx.PublicKey = pvk.PublicKey()
	assert.Nil(t, tx.Verify())
	tx.Data = []byte("fooo")
	assert.NotNil(t, tx.Verify())
}