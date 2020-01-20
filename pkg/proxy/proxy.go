package proxy

import (
	"io"
	"net"

	"github.com/dearcode/crab/log"

	"github.com/dearcode/netpi/pkg/pool"
)

type transport struct {
	remote net.Conn
	agent  net.Conn
}

type Proxy struct {
	token string
}

func New(token string) *Proxy {
	return &Proxy{
		token: token,
	}
}

func (s *Proxy) String() string {
	return s.token
}

func (t *transport) run() {
	go func() {
		io.Copy(t.remote, t.agent)
	}()

	io.Copy(t.agent, t.remote)
}

func newTransport(agent, remote net.Conn) *transport {
	return &transport{
		agent:  agent,
		remote: remote,
	}

}

func (s *Proxy) Run(conn net.Conn) {
	agent := pool.Factory.Wait(s.token, conn)
	log.Infof("%v find agent:%v", s, agent)
	go newTransport(agent, conn).run()
}
