package agent

import (
	"fmt"
	"io"
	"net"

	"github.com/dearcode/crab/log"
	"github.com/dearcode/netpi/pkg/meta"
)

type server struct {
	id        string
	localAddr string
	proxyAddr string
	local     net.Conn
	proxy     net.Conn
}

func newServer(localAddr, proxyAddr, jobID string) *server {
	return &server{
		id:        jobID,
		localAddr: localAddr,
		proxyAddr: proxyAddr,
	}
}

func (s *server) Connect() error {
	conn, err := net.Dial("tcp", s.localAddr)
	if err != nil {
		log.Errorf("connect local:%v error:%v", s.localAddr, err)
		return err
	}
	s.local = conn

	if conn, err = net.Dial("tcp", s.proxyAddr); err != nil {
		s.local.Close()
		log.Errorf("connect proxy:%v error:%v", s.proxyAddr, err)
		return err
	}
	s.proxy = conn

	id := fmt.Sprintf("%*s", meta.MaxTokenLen, s.id)

	_, err = conn.Write([]byte(id))
    if err != nil {
        log.Errorf("write %v, buf:%v, error:%v", s.proxyAddr, id, err)
    }

	return err
}

func (s *server) cleanup() {
	s.local.Close()
	s.proxy.Close()
}

func (s *server) Run() {
	defer s.cleanup()

	go io.Copy(s.local, s.proxy)

	io.Copy(s.proxy, s.local)
}
