package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/yldoge/edu-go-blockchain/crypto"
	"github.com/yldoge/edu-go-blockchain/types"
)

type Header struct {
	Version      uint32
	DataHash     types.Hash
	PreBlockHash types.Hash
	Height       uint32
	Timestamp    int64
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(h); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (h *Header) Hash(hasher Hasher[*Header]) types.Hash {
	return hasher.Hash(h)
}

type Block struct {
	*Header
	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// Cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txs []*Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: txs,
	}, nil
}

func NewBlockFromPrevHeader(prevHeader *Header, txs []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txs)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:      1,
		Height:       prevHeader.Height + 1,
		DataHash:     dataHash,
		PreBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:    time.Now().UnixNano(),
	}

	return NewBlock(header, txs)
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = b.Header.Hash(hasher)
	}
	return b.hash
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Sign(pvk crypto.PrivateKey) error {
	sig, err := pvk.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}
	b.Validator = pvk.PublicKey()
	b.Signature = sig
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}
	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.Hash(BlockHasher{}))
	}

	return nil
}

func CalculateDataHash(txs []*Transaction) (hash types.Hash, err error) {
	buf := &bytes.Buffer{}

	for _, tx := range txs {
		if err = tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return
		}
	}

	hash = sha256.Sum256(buf.Bytes())

	return
}
