package protocol

import (
//"errors"
)

func Parse(src []byte, key *CipherKey) (packet *Packet, cmd Command, err error) {
	packet, err = ExpoundPacket(src, key)
	if err != nil {
		return
	}
	cmd, err = ExpoundCommand(packet)
	return
}
