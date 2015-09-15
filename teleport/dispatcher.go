package teleport

import (
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
	"net"
	"runtime/debug"
)

func Dispatch(data []byte, raddr net.Addr) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(string(debug.Stack()))
		}
	}()

	packet, cmd, err := protocol.Parse(data)

	if err != nil {
		log.Info("wrong data from teleport")
		return
	}

	switch cmd.GetOp() {
	case "1":
	default:
		log.Info("no handler for this command")
	}
	Post2Rails(packet, raddr)
}
