package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yldoge/edu-go-blockchain/crypto"
	"github.com/yldoge/edu-go-blockchain/types"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:      1,
		PreBlockHash: types.RandomHash(),
		Height:       height,
		Timestamp:    time.Now().UnixNano(),
	}
	tx := Transaction{
		Data: []byte("foo"),
	}

	return NewBlock(header, []Transaction{tx})
}

func TestHashBlock(t *testing.T) {
	b := randomBlock(0)
	b.Hash(BlockHasher{})
	assert.NotNil(t, b.hash)
}

func TestSignBlock(t *testing.T) {
	pvk := crypto.GeneratePrivateKey()
	b := randomBlock(0)

	assert.Nil(t, b.Sign(pvk))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	pvk := crypto.GeneratePrivateKey()
	b := randomBlock(0)

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

func TestBlockHeaderDataDiff(t *testing.T) {
	b := randomBlock(0)
	hd := b.HeaderData()
	b.Height = 10
	hd2 := b.HeaderData()
	assert.NotEqual(t, hd, hd2)
}
