package teleport

import (
	"errors"
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
)

type Pool struct {
	conns             map[string]*connection
	unauthorizedConns map[string]*connection
	addrs             map[int64]string
}

// initial vector
func (p *Pool) set_iv(uuid string, iv []byte) {
	c, ok := p.unauthorizedConns[uuid]
	if ok {
		ck := c.cipher_key
		ck.IV = "not used field"
		ck.Iv96str = iv[:12]
		ck.EncryptCtr = 0
		ck.DecryptCtr = 0
	}
}

func (p *Pool) set_user_key(uuid string, uk []byte) {
	c, ok := p.unauthorizedConns[uuid]
	if ok {
		c.cipher_key.UserKey = uk
	}
}

func (p *Pool) set_user_key_index(uuid string, index int) {
	c, ok := p.unauthorizedConns[uuid]
	if ok {
		c.cipher_key.UserKeyIndex = index
	}
}

func (p *Pool) set_teleport_addr(uuid string, addr int64) {
	p.addrs[addr] = uuid
}

func (p *Pool) send(uuid string, packet *protocol.PacketSend, ck *protocol.CipherKey) {

	log.WithFields(log.Fields{
		"uuid":   uuid,
		"packet": packet,
	}).Info("send to teleport")

	var ok bool
	var c *connection
	// teleport login step
	if packet.Op == 2 || packet.Op == 4 {
		c, ok = p.unauthorizedConns[uuid]
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

func (p *Pool) disconnect(key string) {
	c, ok := p.unauthorizedConns[key]
	if ok {
		log.WithFields(log.Fields{
			"remote addr": c.conn.RemoteAddr().String(),
		}).Info("auto close unauthorized connection")
		c.conn.Close()
		delete(p.unauthorizedConns, key)
	}
}

func (p *Pool) get_cipher_key(uuid string) (*protocol.CipherKey, error) {
	c := GlobalPool.conns[uuid]
	if c == nil {
		c = GlobalPool.unauthorizedConns[uuid]
	}

	if c == nil {
		return &protocol.CipherKey{}, errors.New("unknow this uuid: " + uuid)
	}
	return c.cipher_key, nil
}
