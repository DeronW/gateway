package teleport

import (
	"gateway/configs"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom_tcp"
)

func Run() {

	config := &tcp.ServerConfig{
		Host:       configs.TCP_HOST,
		Port:       configs.TCP_PORT,
		Net:        "tcp",
		SendBuf:    configs.TCP_SEND_BUFFER,
		ReceiveBuf: configs.TCP_RECEIVE_BUFFER,
		Deadline:   configs.TCP_DEADLINE,
	}

	handler := tcp.Handler(&Pool{})

	server := tcp.NewServer(config, handler)
	log.Info("start tcp server")
	server.Stop()
	log.Info("stop tcp server")
}
