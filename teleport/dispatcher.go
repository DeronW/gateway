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

	key, err := GetCipherKey(uuid)
	if err != nil {
		return
	}
	packet, cmd, err := protocol.Parse(data, key)

	if err != nil {
		return
	}

	//log.Info("---------------")
	//log.Info(packet)
	//log.Info(cmd)

	switch cmd.GetOp() {
	case "1":
		Post2Rails(packet, uuid)
	case "2":
		//
	default:
		log.Info("no handler for this command")
	}
	return nil
}
