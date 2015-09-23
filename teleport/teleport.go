package teleport

import (
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
)

var GlobalPool *Pool

func Run(cfg *tcp.ServerConfig) {

	GlobalPool = &Pool{
		conns:             make(map[string]*connection, 1000),
		unauthorizedConns: make(map[string]*connection, 1000),
	}

	handler := tcp.Handler(GlobalPool)

	server := tcp.NewServer(cfg, handler)
	log.Info("starting tcp server")
	log.Fatal(server.Start())
	//server.Stop()
	//log.Info("stop tcp server")
}
