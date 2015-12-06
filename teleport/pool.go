package teleport

import (
	"errors"
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
)

type Pool struct {
	conns              map[string]*connection
	unauthorized_conns map[string]*connection
	addrs              map[int]string
}

func (p *Pool) authorize(uuid string) {
	c, ok := p.unauthorized_conns[uuid]
	if ok {
		delete(p.conns, uuid)
		p.conns[uuid] = c
		delete(p.unauthorized_conns, uuid)
	}
}

// initial vector
func (p *Pool) set_iv(uuid string, iv []byte) {
	c, ok := p.unauthorized_conns[uuid]
	if ok {
		ck := c.cipher_key
		ck.IV = "not used field"
		ck.Iv96str = iv[:12]
		ck.EncryptCtr = 0
		ck.DecryptCtr = 0
	}
}

func (p *Pool) set_user_key(uuid string, uk []byte) {
	c, ok := p.unauthorized_conns[uuid]
	if ok {
		c.cipher_key.UserKey = uk
	}
}

func (p *Pool) set_user_key_index(uuid string, index int) {
	c, ok := p.unauthorized_conns[uuid]
	if ok {
		c.cipher_key.UserKeyIndex = index
	}
}

func (p *Pool) set_teleport_addr(uuid string, addr int) {
	p.addrs[addr] = uuid
	c, _ := p.unauthorized_conns[uuid]
	c.addr = addr
}

func (p *Pool) get_teleport_addr(uuid string) (int, bool) {
	c, ok := p.conns[uuid]
	if ok {
		return c.addr, true
	}
	return 0, false
}

func (p *Pool) send(uuid string, packet *protocol.PacketSend, ck *protocol.CipherKey) {

	log.WithFields(log.Fields{
		"uuid":   uuid,
		"packet": packet,
	}).Info("send to teleport")

	var ok bool
	var c *connection
	// teleport login step
	if packet.Op == 2 { // || packet.Op == 4 {
		c, ok = p.unauthorized_conns[uuid]
	} else {
		c, ok = p.conns[uuid]
	}

	if ok {
		err := c.send(packet, ck)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Info("send to teleport error")
		}
	} else {
		log.WithFields(log.Fields{
			"uuid": uuid,
		}).Info("can not find connection")
	}
}

func (p *Pool) disconnect(uuid string) {
	c, ok := p.unauthorized_conns[uuid]
	if ok {
		log.WithFields(log.Fields{
			"remote addr": c.conn.RemoteAddr().String(),
		}).Info("auto close unauthorized connection")
		c.conn.Close()
		delete(p.unauthorized_conns, uuid)
	}

	c2, ok := p.conns[uuid]
	if ok {
		log.WithFields(log.Fields{
			"remote addr": c2.conn.RemoteAddr().String(),
		}).Info("auto close authorized connection")
		c2.conn.Close()
		delete(p.conns, uuid)
	}
}

func (p *Pool) get_cipher_key(uuid string) (*protocol.CipherKey, error) {
	c := GlobalPool.conns[uuid]
	if c == nil {
		c = GlobalPool.unauthorized_conns[uuid]
	}

	if c == nil {
		return &protocol.CipherKey{}, errors.New("unknow this uuid: " + uuid)
	}
	return c.cipher_key, nil
}
