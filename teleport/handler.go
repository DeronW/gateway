package teleport

import (
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
)

type Pool struct {
	conns map[int]*tcp.Conn
}

func (p *Pool) OnConnect(c *tcp.Conn) bool {
	log.WithFields(log.Fields{
		"remote addr": c.RemoteAddr(),
	}).Info("connected")
	return true
}

func (p *Pool) OnMessage(c *tcp.Conn, m []byte) bool {
	log.WithFields(log.Fields{"message": m}).Info("receive")
	return true
}

func (p *Pool) OnClose(c *tcp.Conn) {
	log.WithFields(log.Fields{
		"addr": c.RemoteAddr(),
	}).Info("closed")
}
