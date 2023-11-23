package core

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (bv *BlockValidator) ValidateBlock(b *Block) error {
	if bv.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%d)", b.Height, b.Hash(BlockHasher{}))
	}

	if b.Height > bv.bc.Height()+1 {
		return fmt.Errorf("block (%s) too high", b.Hash(BlockHasher{}))
	}

	prevHeader, err := bv.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}

	hash := prevHeader.Hash(BlockHasher{})
	if hash != b.PreBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PreBlockHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
