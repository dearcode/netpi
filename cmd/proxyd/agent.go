package main

import (
	"net"
	"strings"

	"github.com/dearcode/crab/log"

	"github.com/dearcode/netpi/pkg/meta"
	"github.com/dearcode/netpi/pkg/pool"
)

type agentServer struct {
	listener net.Listener
}

func newAgentServer(addr string) *agentServer {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("start listen:%v, error:%v", addr, err)
	}

	return &agentServer{
		listener: listener,
	}
}

func (s *agentServer) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Errorf("%v listener Accept error:%v", s, err)
			return
		}

		buf := make([]byte, meta.MaxTokenLen)
		if _, err = conn.Read(buf); err != nil {
			log.Errorf("%v read agent:%v error:%v", s, conn, err)
			conn.Close()
			continue
		}

		id := strings.TrimSpace(string(buf))
		pool.Factory.PutWorker(id, conn)
		log.Infof("%v new worker:%v from:%v", s, id, conn.RemoteAddr())
	}

}
