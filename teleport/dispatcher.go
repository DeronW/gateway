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
	packet, cmd, err := protocol.Parse(data, ck)

	if err != nil {
		return
	}

	switch cmd.GetOp() {
	case "1":
		pack4send, iv, uk, uki := protocol.LoginStepOne(packet, ck)

		GlobalPool.set_iv(uuid, iv)
		GlobalPool.set_user_key(uuid, uk)
		GlobalPool.set_user_key_index(uuid, uki)

		GlobalPool.send(uuid, pack4send, ck)
	case "3":
		GlobalPool.send(uuid, &protocol.PacketSend{
			Encrypted:         true,
			WirelessEncrypted: true,
			Op:                4,
			Params:            "",
			Version:           packet.Version,
		}, ck)
	case "qt":
		log.Info("return time")
	default:
		log.Info("no handler for this command")
	}
	return nil
}
