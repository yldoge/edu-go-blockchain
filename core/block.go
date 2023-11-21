package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

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

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// Cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txs []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
	}
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}
	return b.hash
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(b.Header); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (b *Block) Sign(pvk crypto.PrivateKey) error {
	sig, err := pvk.Sign(b.HeaderData())
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
	if !b.Signature.Verify(b.Validator, b.HeaderData()) {
		return fmt.Errorf("block has invalid signature")
	}
	return nil
}
