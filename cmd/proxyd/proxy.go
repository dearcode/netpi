package main

import (
	"net"

	"dearcode.net/crab/log"

	"dearcode.net/netpi/pkg/proxy"
)

type proxyServer struct {
	token    string
	listener net.Listener
	server   *proxy.Proxy
}

func newProxyServer(addr, token string) *proxyServer {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("start listen:%v, error:%v", addr, err)
	}

	return &proxyServer{
		listener: listener,
		token:    token,
		server:   proxy.New(token),
	}
}

func (s *proxyServer) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Errorf("%v listener Accept error:%v", s, err)
			return
		}

		s.server.Run(conn)

		log.Infof("%v new conn from:%v", s, conn.RemoteAddr())
	}

}
