package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack(128)

	s.Push(1)
	s.Push(2)

	assert.Equal(t, 1, s.Pop())
	assert.Equal(t, 2, s.Pop())
}

func TestPushIntThenAdd(t *testing.T) {
	data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0b}
	vm := NewVM(data, NewState())
	assert.Nil(t, vm.Run())

	assert.Equal(t, 3, vm.stack.Pop())
}

func TestPushIntThenSubtract(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}
	vm := NewVM(data, NewState())
	assert.Nil(t, vm.Run())

	assert.Equal(t, 1, vm.stack.Pop())

	data = []byte{0x01, 0x0a, 0x02, 0x0a, 0x0e}
	vm = NewVM(data, NewState())
	assert.Nil(t, vm.Run())

	assert.Equal(t, -1, vm.stack.Pop())
}

func TestPushByteThenPack(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d}
	vm := NewVM(data, NewState())
	assert.Nil(t, vm.Run())

	assert.Equal(t, "FOO", string(vm.stack.Pop().([]byte)))
}

func TestContractStore(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	valueBytes, err := contractState.Get([]byte("FOO"))
	assert.Nil(t, err)
	value := deserializeInt64(valueBytes)
	assert.Equal(t, value, int64(5))
}