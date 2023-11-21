package core

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yldoge/edu-go-blockchain/types"
)

var (
	h *Header = &Header{
		Version:   1,
		PreBlock:  types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:    10,
		Nonce:     989394,
	}
	b *Block = &Block{
		Header:       *h,
		Transactions: nil,
	}
)

func TestBlock_Header_Encode_Decode(t *testing.T) {

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}

func TestBlock_Encode_Decode(t *testing.T) {
	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, b, bDecode)
}

func TestBlock_Hash(t *testing.T) {
	var headerHash types.Hash
	headerHash = b.hash
	assert.True(t, headerHash.IsZero())
	headerHash = b.Hash()
	fmt.Println(headerHash)
	assert.False(t, headerHash.IsZero())

	// reset global value
	b.hash = types.HashFromBytes(make([]byte, 32))
}
