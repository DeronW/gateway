package teleport

import (
	"gateway/config"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	"sync"
)

var tele_once sync.Once
var GlobalPool *Pool

var TCP_CONFIG *config.TCP
var KEEPALIVE_CONFIG *config.KeepAlive

func get_teleport_config() *tcp.ServerConfig {
	return &tcp.ServerConfig{
		Host:       TCP_CONFIG.Host,
		Port:       uint32(TCP_CONFIG.Port),
		Net:        "tcp",
		SendBuf:    uint32(TCP_CONFIG.SendBuffer),
		ReceiveBuf: uint32(TCP_CONFIG.ReceiveBuffer),
		Separtor:   TCP_CONFIG.Separtor,

		Deadline:          KEEPALIVE_CONFIG.Deadline,
		KeepAlive:         KEEPALIVE_CONFIG.Enable,
		KeepAliveIdle:     KEEPALIVE_CONFIG.Idle,
		KeepAliveCount:    KEEPALIVE_CONFIG.Count,
		KeepAliveInterval: KEEPALIVE_CONFIG.Interval,
	}
}

func Run() {
	GlobalPool = &Pool{
		conns:              make(map[string]*connection, 1000),
		unauthorized_conns: make(map[string]*connection, 1000),
		addrs:              make(map[int]string, 1000),
	}

	handler := tcp.Handler(GlobalPool)

	server := tcp.NewServer(get_teleport_config(), handler)
	log.Info("starting tcp server")
	log.Fatal(server.Start())
	//server.Stop()
	//log.Info("stop tcp server")
}

func init() {
	tele_once.Do(func() {
		TCP_CONFIG = config.GetTCPConfig()
		KEEPALIVE_CONFIG = config.GetKeepAliveCfg()
	})
}
