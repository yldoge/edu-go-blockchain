package network

import (
	"sync"

	"github.com/yldoge/edu-go-blockchain/core"
	"github.com/yldoge/edu-go-blockchain/types"
)

type TxPool struct {
	all     *TxSortedMap
	pending *TxSortedMap
	// When the pool is full, prune the oldest transaction
	maxLen int
}

func NewTxPool(maxLen int) *TxPool {
	return &TxPool{
		all:     NewTxSortedMap(),
		pending: NewTxSortedMap(),
		maxLen:  maxLen,
	}
}

func (p *TxPool) Add(tx *core.Transaction) {
	if p.all.Count() == p.maxLen {
		oldest := p.all.First()
		p.all.Remove(oldest.Hash(core.TxHasher{}))
	}

	if !p.all.Contains(tx.Hash(core.TxHasher{})) {
		p.all.Add(tx)
		p.pending.Add(tx)
	}
}

func (p *TxPool) Contains(hash types.Hash) bool {
	return p.all.Contains(hash)
}

// Pending returns a slice of transactions that are in the pending pool
func (p *TxPool) Pending() []*core.Transaction {
	return p.pending.txs.Data
}

func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}

type TxSortedMap struct {
	lock   sync.RWMutex
	lookup map[types.Hash]*core.Transaction
	txs    *types.List[*core.Transaction]
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(map[types.Hash]*core.Transaction),
		txs:    types.NewList[*core.Transaction](),
	}
}

func (t *TxSortedMap) First() *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	first := t.txs.Get(0)
	return t.lookup[first.Hash(core.TxHasher{})]
}

func (t *TxSortedMap) Get(h types.Hash) *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.lookup[h]
}

func (t *TxSortedMap) Add(tx *core.Transaction) {
	hash := tx.Hash(core.TxHasher{})

	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.lookup[hash]; !ok {
		t.lookup[hash] = tx
		t.txs.Insert(tx)
	}
}

func (t *TxSortedMap) Remove(h types.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.txs.Remove(t.lookup[h])
	delete(t.lookup, h)
}

func (t *TxSortedMap) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.lookup)
}

func (t *TxSortedMap) Contains(h types.Hash) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	_, ok := t.lookup[h]
	return ok
}

func (t *TxSortedMap) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.lookup = make(map[types.Hash]*core.Transaction)
	t.txs.Clear()
}
