package teleport

import (
	"gateway/configs"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
)

func Run() {

	config := &tcp.ServerConfig{
		Host:       configs.TCP_HOST,
		Port:       configs.TCP_PORT,
		Net:        "tcp",
		SendBuf:    configs.TCP_SEND_BUFFER,
		ReceiveBuf: configs.TCP_RECEIVE_BUFFER,

		Deadline:          configs.TCP_DEADLINE,
		KeepAlive:         configs.TCP_KEEP_ALIVE,
		KeepAliveIdle:     configs.TCP_KEEP_ALIVE_IDLE,
		KeepAliveCount:    configs.TCP_KEEP_ALIVE_COUNT,
		KeepAliveInterval: configs.TCP_KEEP_ALIVE_INTERVAL,

		Separtor: configs.TCP_SERARTOR,
	}

	handler := tcp.Handler(&Pool{})

	server := tcp.NewServer(config, handler)
	log.Info("starting tcp server")
	log.Fatal(server.Start())
	//server.Stop()
	//log.Info("stop tcp server")
}
