package main

import (
	"flag"
	"github.com/dearcode/netpi/pkg/agent"
)

var (
	local = flag.String("local", "127.0.0.1:8081", "local server addr")
	proxy = flag.String("proxy", "127.0.0.1", "proxy server addr, to fetch job list, connect proxy")
	id    = flag.String("id", "12345678", "worker id")
)

func main() {
	flag.Parse()

	m := agent.NewManager(*id, *local, *proxy)
    m.Run()
}
