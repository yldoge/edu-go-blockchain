package network

import (
	"io"
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
	b, err := io.ReadAll(rpc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)
	assert.Equal(t, rpc.From, trl.addr)
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("foo")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <-trb.Consume()
	b, err := io.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcc := <-trc.Consume()
	c, err := io.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, c, msg)
}
