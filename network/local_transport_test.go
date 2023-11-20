package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	trl := NewLocalTransport("LOCAL")
	trr := NewLocalTransport("REMOTE")

	trl.Connect(trr)
	trr.Connect(trl)
	assert.Equal(t, trl.peers[trr.addr], trr)
	assert.Equal(t, trr.peers[trl.addr], trl)
}

func TestSendMessage(t *testing.T) {
	trl := NewLocalTransport("LOCAL")
	trr := NewLocalTransport("REMOTE")

	trl.Connect(trr)
	trr.Connect(trl)

	msg := []byte("Hello, world")
	assert.Nil(t, trl.SendMessage(trr.addr, msg))

	rpc := <-trr.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, trl.addr)
}
