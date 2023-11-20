package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"io"

	"github.com/yldoge/edu-go-blockchain/types"
)

type Header struct {
	Version uint32

	PreBlock  types.Hash
	Timestamp int64
	Height    uint32
	Nonce     uint64
}

func (h *Header) EncodeBinary(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(h); err != nil {
		return err
	}
	return nil
}

func (h *Header) DecodeBinary(r io.Reader) error {
	dec := gob.NewDecoder(r)
	if err := dec.Decode(h); err != nil {
		return err
	}
	return nil
}

type Block struct {
	Header
	Transactions []Transaction

	// Cached version of the header hash
	hash types.Hash
}

func (b *Block) Hash() types.Hash {
	if !b.hash.IsZero() {
		return b.hash
	}
	buf := &bytes.Buffer{}
	b.Header.EncodeBinary(buf)
	b.hash = types.Hash(sha256.Sum256(buf.Bytes()))

	return b.hash
}

func (b *Block) EncodeBinary(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(b); err != nil {
		return err
	}
	return nil
}

func (b *Block) DecodeBinary(r io.Reader) error {
	dec := gob.NewDecoder(r)
	if err := dec.Decode(b); err != nil {
		return err
	}
	return nil
}
