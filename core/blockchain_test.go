package core

import (
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/yldoge/edu-go-blockchain/types"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(log.NewLogfmtLogger(os.Stderr), randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)

	return bc
}

func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(100))
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	chainHeight := 1000
	for i := 0; i < chainHeight; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, len(bc.headers), 1001)
	assert.Equal(t, bc.Height(), uint32(1000))
	// Add a unsigned block should fail
	assert.NotNil(t, bc.AddBlock(randomBlock(t, uint32(chainHeight)+1, types.Hash{})))
	// Add an existed block should fail
	assert.NotNil(t, bc.AddBlock(randomBlock(t, uint32(100), types.Hash{})))
}

func TestAddBlockToHeight(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlock(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	chainHeight := 100
	for i := 0; i < chainHeight; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint32(i + 1))

		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, h uint32) types.Hash {
	prevHeader, err := bc.GetHeader(h - 1)
	assert.Nil(t, err)
	return prevHeader.Hash(BlockHasher{})
}
