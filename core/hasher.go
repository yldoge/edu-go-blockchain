package core

import (
	"crypto/sha256"

	"github.com/yldoge/edu-go-blockchain/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(h *Header) types.Hash {
	shaH := sha256.Sum256(h.Bytes())
	return types.Hash(shaH)
}

type TxHasher struct{}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data))
}
