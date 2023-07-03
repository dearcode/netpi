package main

import (
	"flag"

	"dearcode.net/crab/log"
	"dearcode.net/doodle/pkg/service"

	_ "dearcode.net/netpi/pkg/pool"
	"dearcode.net/netpi/pkg/proxy"
)

type server struct {
}

var (
	s         server
	debug     = flag.Bool("debug", false, "debug")
	agentAddr = flag.String("agent", ":9877", "agent listen address")
	proxyAddr = flag.String("proxy", ":9878", "proxy listen address")
	token     = flag.String("token", "12345678", "token")
)

func main() {
	flag.Parse()

	if !*debug {
		log.SetLevel(log.LogInfo)
		log.SetColor(false)
		log.SetOutputFile("./logs/server.log")
	}

	p := newProxyServer(*proxyAddr, *token)
	go p.Run()

	a := newAgentServer(*agentAddr)
	go a.Run()

	srv := service.New()
	srv.Init()
	srv.Register(proxy.Job{})
	srv.Start()
}
