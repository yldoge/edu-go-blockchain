package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/yldoge/edu-go-blockchain/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{key: key}
}

func (pvk PrivateKey) PublicKey() PublicKey {
	return PublicKey{Key: &pvk.key.PublicKey}
}

func (pvk PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, pvk.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (pbk PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(pbk.Key, pbk.Key.X, pbk.Key.Y)
}

func (pbk PublicKey) Address() types.Address {
	h := sha256.Sum256(pbk.ToSlice())

	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	S, R *big.Int
}

func (sig Signature) Verify(pbk PublicKey, data []byte) bool {
	return ecdsa.Verify(pbk.Key, data, sig.R, sig.S)
}
