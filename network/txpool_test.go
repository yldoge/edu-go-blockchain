package network

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yldoge/edu-go-blockchain/core"
)

func TextTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
}

func TestTxPoolAdd(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("foo"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)

	txx := core.NewTransaction([]byte("foo"))
	assert.Nil(t, p.Add(txx))
	assert.Equal(t, p.Len(), 1)

	p.Flush()
	assert.Equal(t, p.Len(), 0)
}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000

	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(rand.Int63n(500000))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, p.Len(), txLen)

	txs := p.Transactions()
	for i := 0; i < txLen-1; i++ {
		assert.True(t, txs[i].FirstSeen() <= txs[i+1].FirstSeen())
	}
}
