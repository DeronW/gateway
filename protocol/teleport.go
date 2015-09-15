package protocol

import (
	"gateway/protocol/command"
)

func Parse(src []byte) (packet *command.Packet, cmd command.Command, err error) {
	packet, err = command.ExpoundPacket(src)
	if err != nil {
		return
	}
	cmd, err = command.ExpoundCommand(packet)
	return
}
