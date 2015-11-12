package teleport

import (
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
	"runtime/debug"
)

func Dispatch(data []byte, uuid string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(string(debug.Stack()))
		}
	}()

	ckey, err := GetCipherKey(uuid)
	if err != nil {
		return
	}
	packet, cmd, err := protocol.Parse(data, ckey)

	if err != nil {
		return
	}

	switch cmd.GetOp() {
	case "1", "3":
		Post2RailsLoginCmd(packet, uuid, ckey)
	case "2", "4":
		//
	default:
		log.Info("no handler for this command")
	}
	return nil
}
