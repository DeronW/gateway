package teleport

import (
	"errors"
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	"time"
)

type connection struct {
	conn         *tcp.Conn
	iv           string
	iv96str      string
	encryptCtr   int32
	decryptCtr   int32
	userKey      string
	userKeyIndex int
	authorized   bool
}

type Pool struct {
	conns             map[string]*connection
	unauthorizedConns map[string]*connection
}

// SetIv set 4 default value
func (p *Pool) SetIV(uuid string, iv string) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		c.iv = iv
		// this is actualy copy of protocol/tools.go: func reverse
		a := []byte(iv)
		t := len(a) - 1
		b := make([]byte, t+1)
		for i := 0; i <= t; i++ {
			b[i] = a[t-i]
		}
		// copy end
		c.iv96str = string(b)
		c.encryptCtr = 0
		c.decryptCtr = 0
	}
}

func (p *Pool) SetUserKey(uuid string, uk string) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		c.userKey = uk
	}
}

func (p *Pool) SetUserKeyIndex(uuid string, index int) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		c.userKeyIndex = index
	}
}

func GetCipherKey(uuid string) (*protocol.CipherKey, error) {
	c := GlobalPool.conns[uuid]
	if c == nil {
		c = GlobalPool.unauthorizedConns[uuid]
	}

	if c == nil {
		return &protocol.CipherKey{}, errors.New("unknow this uuid: " + uuid)
	}

	return &protocol.CipherKey{
		UserKeyIndex: c.userKeyIndex,
		IV:           c.iv,
		Iv96str:      c.iv96str,
		EncryptCtr:   c.encryptCtr,
		DecryptCtr:   c.decryptCtr,
		UserKey:      c.userKey,
	}, nil
}

func (p *Pool) Send(uuid string, msg string) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	log.WithFields(log.Fields{
		"message": msg,
	}).Info("send to teleport")

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
	p.unauthorizedConns[c.RemoteAddr().String()] = &connection{conn: c, authorized: false}
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
	err := Dispatch(m, c.RemoteAddr().String())
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
