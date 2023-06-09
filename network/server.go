package network

import (
	"fmt"
	"time"

	"github.com/peyzor/xchain/core"
	"github.com/peyzor/xchain/crypto"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
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
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	return &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}

	fmt.Println("server shutdown")
}

func (s *Server) handleTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	if err := s.memPool.Add(tx); err != nil {
		logrus.Infoln(err)
		return err
	}

	logrus.WithFields(logrus.Fields{
		"hash": tx.Hash(core.TxHasher{}),
	}).Info("adding new tx to the mempool")
	return nil
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
}

func (s *Server) initTransport() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)

	}
}
