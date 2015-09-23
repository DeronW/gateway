package teleport

import (
	"gateway/configs"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	"time"
)

type connection struct {
	conn         *tcp.Conn
	iv           string
	userKey      string
	userKeyIndex int
}

type Pool struct {
	conns             map[string]*connection
	unauthorizedConns map[string]*connection
}

func (p *Pool) SetIV(uuid string, iv string) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	if ok {
		c.iv = iv
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

func (p *Pool) Send(uuid string, msg string) {
	c, ok := GlobalPool.unauthorizedConns[uuid]
	log.Info("===================== send =====================  ")
	log.Info(msg)
	if ok {
		err := c.conn.Write([]byte(msg), 0)
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
	timeout := time.Second * configs.TCP_AUTO_CLOSE_DURATION
	if timeout != 0 {
		time.AfterFunc(timeout, func() {
			if c.Id == 0 {
				p.disconnect(c.RemoteAddr().String())
			}
		})
	}
	p.unauthorizedConns[c.RemoteAddr().String()] = &connection{conn: c}
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
	Dispatch(m, c.RemoteAddr().String())
	return true
}

func (p *Pool) OnClose(c *tcp.Conn) {
	log.WithFields(log.Fields{
		"remote addr": c.RemoteAddr(),
	}).Info("closed")
}
