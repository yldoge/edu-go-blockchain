package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0))
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
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	chainHeight := 1000
	for i := 0; i < chainHeight; i++ {
		block := randomBlockWithSignature(t, uint32(i+1))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, len(bc.headers), 1001)
	assert.Equal(t, bc.Height(), uint32(1000))
	// Add a unsigned block should fail
	assert.NotNil(t, bc.AddBlock(randomBlock(uint32(chainHeight)+1)))
	// Add an existed block should fail
	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, uint32(100))))
}
