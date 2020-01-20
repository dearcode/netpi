package main

import (
	"flag"

	"github.com/dearcode/crab/log"

	"github.com/dearcode/netpi/pkg/agent"
)

var (
	debug = flag.Bool("debug", false, "debug")
	local = flag.String("local", "127.0.0.1:8081", "local server addr")
	proxy = flag.String("proxy", "127.0.0.1", "proxy server addr, to fetch job list, connect proxy")
	id    = flag.String("id", "12345678", "worker id")
)

func main() {
	flag.Parse()

	if !*debug {
		log.SetLevel(log.LogInfo)
		log.SetColor(false)
		log.SetOutputFile("./logs/server.log")
	}

	m := agent.NewManager(*id, *local, *proxy)
	m.Run()
}
