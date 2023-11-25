package network

import (
	"bytes"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yldoge/edu-go-blockchain/core"
	"github.com/yldoge/edu-go-blockchain/crypto"
)

var (
	defaultBlockTime = 5 * time.Second
)

type ServerOpts struct {
	RPCDecodeFunc
	RPCProcessor
	Transports []Transport
	BlockTime  time.Duration
	PrivateKey *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == 0 {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	// server itself is the default rpc processor
	if opts.RPCProcessor == nil {
		opts.RPCProcessor = s
	}

	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				log.Error(err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				log.Error(err)
			}
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNetwork()
			}
		}
	}
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

// Implementing RPCProcessor and routing to correct process func
func (s *Server) ProcessMessage(decMsg *DecodedMessage) error {
	switch t := decMsg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {

	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		log.WithFields(log.Fields{
			"hash": hash,
		}).Info("transaction already in mempool")
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	log.WithFields(log.Fields{
		"hash":       hash,
		"mempoolLen": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) createNetwork() error {
	fmt.Println("creating a new block")
	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
