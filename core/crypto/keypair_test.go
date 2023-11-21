package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifySignatureSuccess(t *testing.T) {
	pvk := GeneratePrivateKey()
	pbk := pvk.PublicKey()

	msg := []byte("Hello, world")
	sig, err := pvk.Sign(msg)
	assert.Nil(t, err)

	assert.True(t, sig.Verify(pbk, msg))
}

func TestVerifySignature_Payload_Modified(t *testing.T) {
	pvk := GeneratePrivateKey()
	pbk := pvk.PublicKey()

	msg := []byte("Hello, world")
	sig, err := pvk.Sign(msg)
	assert.Nil(t, err)

	modifiedMsg := []byte("Hello, world!!")
	assert.False(t, sig.Verify(pbk, modifiedMsg))
}

func TestVerifySignature_Invalid_PublicKey(t *testing.T) {
	pvk := GeneratePrivateKey()
	otherPbk := GeneratePrivateKey().PublicKey()

	msg := []byte("Hello, world")
	sig, err := pvk.Sign(msg)
	assert.Nil(t, err)
	assert.False(t, sig.Verify(otherPbk, msg))
}
