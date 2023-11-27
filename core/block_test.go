package core

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yldoge/edu-go-blockchain/crypto"
	"github.com/yldoge/edu-go-blockchain/types"
)

func TestHashBlock(t *testing.T) {
	b := randomBlock(t, 0, types.Hash{})
	b.Hash(BlockHasher{})
	assert.NotNil(t, b.hash)
}

func TestSignBlock(t *testing.T) {
	pvk := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(pvk))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	pvk := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

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

func TestBlockCodec(t *testing.T) {
	b := randomBlock(t, 1, types.Hash{})
	buf := &bytes.Buffer{}
	assert.Nil(t, b.Encode(NewGobBlockEncoder(buf)))

	bDecoded := &Block{}
	assert.Nil(t, bDecoded.Decode(NewGobBlockDecoder(buf)))
	assert.Equal(t, b, bDecoded)
}

func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	pvk := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)

	header := &Header{
		Version:      1,
		PreBlockHash: prevBlockHash,
		Height:       height,
		Timestamp:    time.Now().UnixNano(),
	}

	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	assert.Nil(t, b.Sign(pvk))

	return b
}
