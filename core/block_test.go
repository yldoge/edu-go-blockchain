package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yldoge/edu-go-blockchain/crypto"
	"github.com/yldoge/edu-go-blockchain/types"
)

func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:      1,
		PreBlockHash: prevBlockHash,
		Height:       height,
		Timestamp:    time.Now().UnixNano(),
	}

	return NewBlock(header, []Transaction{})
}

func randomBlockWithSignature(t *testing.T, h uint32, prevBlockHash types.Hash) *Block {
	pvk := crypto.GeneratePrivateKey()
	b := randomBlock(h, prevBlockHash)
	tx := randomTxWithSignature(t)
	b.AddTransaction(tx)
	assert.Nil(t, b.Sign(pvk))

	return b
}

func TestHashBlock(t *testing.T) {
	b := randomBlock(0, types.Hash{})
	b.Hash(BlockHasher{})
	assert.NotNil(t, b.hash)
}

func TestSignBlock(t *testing.T) {
	pvk := crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(pvk))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	pvk := crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(pvk))
	assert.Nil(t, b.Verify())

	// Invalid public key should fail the verification
	b.Validator = crypto.GeneratePrivateKey().PublicKey()
	assert.NotNil(t, b.Verify())
	// Use original public key
	b.Validator = pvk.PublicKey()
	assert.Nil(t, b.Verify())
	// Modify header data should fail the verification
	b.Height = 10
	// assert.NotNil(t, b.Verify())

}
