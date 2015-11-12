package teleport

import (
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	"time"
)

var GlobalPool *Pool
var AutoCloseConnection time.Duration
var RailsPostUrl string

func Run(cfg *tcp.ServerConfig, autoCloseConnection time.Duration, railsPostUrl string) {

	GlobalPool = &Pool{
		conns:             make(map[string]*connection, 1000),
		unauthorizedConns: make(map[string]*connection, 1000),
		addrs:             make(map[int64]string, 1000),
	}

	AutoCloseConnection = autoCloseConnection
	RailsPostUrl = railsPostUrl

	handler := tcp.Handler(GlobalPool)

	server := tcp.NewServer(cfg, handler)
	log.Info("starting tcp server")
	log.Fatal(server.Start())
	//server.Stop()
	//log.Info("stop tcp server")
}
