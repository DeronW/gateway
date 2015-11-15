package teleport

import (
	"errors"
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	"strconv"
	"time"
)

type connection struct {
	conn       *tcp.Conn
	cipher_key *protocol.CipherKey
	authorized bool
}

type Pool struct {
	conns             map[string]*connection
	unauthorizedConns map[string]*connection
	addrs             map[int64]string
}

// SetIv set 4 default value
func (p *Pool) SetIV(uuid string, iv string, iv_chr string) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		ckey := c.cipher_key
		ckey.IV = iv
		for i := 0; i < len(iv_chr) && i < 24; i += 2 {
			b, _ := strconv.ParseInt(iv_chr[i:i+2], 16, 16)
			ckey.Iv96str = append(ckey.Iv96str, byte(b))
		}
		ckey.EncryptCtr = 0
		ckey.DecryptCtr = 0
	}
}

func (p *Pool) SetIV2(uuid string, iv []byte) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		ckey := c.cipher_key
		ckey.IV = "not used field"
		ckey.Iv96str = iv[:12]
		ckey.EncryptCtr = 0
		ckey.DecryptCtr = 0
	}
}

func (p *Pool) SetUserKey(uuid string, uk []byte) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		c.cipher_key.UserKey = uk
	}
}

func (p *Pool) SetUserKeyIndex(uuid string, index int) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		c.cipher_key.UserKeyIndex = index
	}
}

func (p *Pool) SetTeleportAddr(uuid string, addr int64) {
	GlobalPool.addrs[addr] = uuid
}

func GetCipherKey(uuid string) (*protocol.CipherKey, error) {
	c := GlobalPool.conns[uuid]
	if c == nil {
		c = GlobalPool.unauthorizedConns[uuid]
	}

	if c == nil {
		return &protocol.CipherKey{}, errors.New("unknow this uuid: " + uuid)
	}
	return c.cipher_key, nil
}

func (p *Pool) Send(uuid string, msg string) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		timeout := time.Second * 5 //  time.Duration
		err := c.conn.Write([]byte(msg), timeout)
		if err != nil {
			log.Info("send to teleport error")
			log.Info(err)
		}
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

func (p *Pool) OnConnect(c *tcp.Conn) bool {
	log.WithFields(log.Fields{
		"remote addr": c.RemoteAddr(),
	}).Info("connected")

	// auto close conn if it does not auth in 60s
	timeout := time.Second * AutoCloseConnection
	if timeout != 0 {
		time.AfterFunc(timeout, func() {
			if c.Id == 0 {
				p.disconnect(c.RemoteAddr().String())
			}
		})
	}
	p.unauthorizedConns[c.RemoteAddr().String()] = &connection{
		conn:       c,
		cipher_key: &protocol.CipherKey{},
		authorized: false,
	}
	return true
}

func (p *Pool) OnMessage(c *tcp.Conn, m []byte) bool {
	defer func() {
		if r := recover(); r != nil {
			log.WithFields(log.Fields{
				"error": r,
			}).Info("runtime error")
		}
	}()

	// Well!!! HERE, last chr of value m is *, so ignore it
	//err := Dispatch(m, c.RemoteAddr().String())
	err := Dispatch(m[:len(m)-1], c.RemoteAddr().String())
	if err != nil {
		log.Info(err)
	}
	return true
}

func (p *Pool) OnClose(c *tcp.Conn) {
	log.WithFields(log.Fields{
		"remote addr": c.RemoteAddr(),
		"time":        time.Now(),
	}).Info("closed")
}
