package teleport

import (
	"gateway/protocol"
	tcp "github.com/delongw/phantom-tcp"
	"time"
)

type connection struct {
	conn       *tcp.Conn
	cipher_key *protocol.CipherKey
	authorized bool
}

func (c *connection) send(
	p *protocol.PacketSend,
	ck *protocol.CipherKey,
) error {
	var err error
	enc, err := protocol.Encrypt(p, ck)

	if err != nil {
		return err
	}

	err = c.conn.Write([]byte(enc), time.Second*5)
	if err != nil {
		return err
	}
	return nil
}
