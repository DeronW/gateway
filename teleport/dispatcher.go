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

	ck, _ := GlobalPool.get_cipher_key(uuid)

	pk, err := protocol.ExpoundPacketReceive(data, ck)

	log.WithFields(log.Fields{"packet": pk}).Info("receive packet")

	if err != nil {
		log.WithFields(log.Fields{
			"packet": pk,
			"error":  err,
		}).Info("expound packet error")
		return
	}

	var cmd protocol.Command
	switch pk.Op {
	case "1":
		c, iv, uk, uki := protocol.CommandLoginSetup(pk, ck)
		cmd = c

		GlobalPool.set_iv(uuid, iv)
		GlobalPool.set_user_key(uuid, uk)
		GlobalPool.set_user_key_index(uuid, uki)

		//GlobalPool.send(uuid, pack4send, ck)
	case "3":
		cmd = &protocol.Command_login3{protocol.CommandBase{pk}}
	case "qt":
		log.Info("return time")
	default:
		log.Info("no handler for this command")
	}

	if cmd != nil {
		handle_command(uuid, cmd, ck)
	}
	return nil
}

func handle_command(uuid string, cmd protocol.Command, ck *protocol.CipherKey) {
	pk, ok := cmd.GetReply()
	if ok {
		GlobalPool.send(uuid, pk, ck)
	}

	msg, ok := cmd.GetMessage()
	if ok {
		log.Info(msg)
	}

	ev, ok := cmd.GetEvent()
	if ok {
		log.Info(ev)
	}
}
