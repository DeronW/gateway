package teleport

import (
	"gateway/configs"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	"time"
)

var GlobalPool *Pool

func Run() {

	config := &tcp.ServerConfig{
		Host:       configs.TCP_HOST,
		Port:       configs.TCP_PORT,
		Net:        "tcp",
		SendBuf:    configs.TCP_SEND_BUFFER,
		ReceiveBuf: configs.TCP_RECEIVE_BUFFER,

		Deadline:          time.Second * configs.TCP_DEADLINE,
		KeepAlive:         configs.TCP_KEEP_ALIVE,
		KeepAliveIdle:     time.Second * configs.TCP_KEEP_ALIVE_IDLE,
		KeepAliveCount:    configs.TCP_KEEP_ALIVE_COUNT,
		KeepAliveInterval: time.Second * configs.TCP_KEEP_ALIVE_INTERVAL,

		Separtor: configs.TCP_SERARTOR,
	}

	GlobalPool = &Pool{
		conns:             make(map[string]*connection, 1000),
		unauthorizedConns: make(map[string]*connection, 1000),
	}

	handler := tcp.Handler(GlobalPool)

	server := tcp.NewServer(config, handler)
	log.Info("starting tcp server")
	log.Fatal(server.Start())
	//server.Stop()
	//log.Info("stop tcp server")
}
