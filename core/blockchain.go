package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
)

type Blockchain struct {
	logger    log.Logger
	lock      sync.RWMutex
	headers   []*Header
	store     Storage
	validator Validator
}

func NewBlockchain(l log.Logger, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStorage(),
		logger:  l,
	}

	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		vm := NewVM(tx.Data, NewState())
		if err := vm.Run(); err != nil {
			return err
		}

		fmt.Printf("STATE: %+v\n", vm.contractState)
		fmt.Printf("VM RESULT: %+v\n", vm.stack.Pop())
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) HasBlock(h uint32) bool {
	return h <= bc.Height()
}

func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) GetHeader(h uint32) (*Header, error) {
	if h > bc.Height() {
		return nil, fmt.Errorf("given height (%d) is too high", h)
	}

	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.headers[h], nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", bc.Height(),
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}
