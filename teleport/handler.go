package teleport

import (
	"gateway/configs"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	"time"
)

type Pool struct {
	conns map[int]*tcp.Conn
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
				c.Close()
			}
		})
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
	Dispatch(m, c.RemoteAddr())
	return true
}

func (p *Pool) OnClose(c *tcp.Conn) {
	log.WithFields(log.Fields{
		"remote addr": c.RemoteAddr(),
	}).Info("closed")
}
