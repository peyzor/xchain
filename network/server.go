package network

import (
	"fmt"
	"time"

	"github.com/peyzor/xchain/crypto"
)

type ServerOpts struct {
	Transports []Transport
	PrivateKey *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts:  opts,
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			fmt.Println("do stuff every x seconds")
		}
	}

	fmt.Println("server shutdown")
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
