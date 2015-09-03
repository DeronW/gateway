package teleport

import (
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	//"time"
	"encoding/base64"
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
	log.WithFields(log.Fields{
		"message": string(m),
	}).Info("receive")

	s, _ := decode(m)

	println(string(s))
	return true
}

func (p *Pool) OnClose(c *tcp.Conn) {
	log.WithFields(log.Fields{
		"addr": c.RemoteAddr(),
	}).Info("closed")
}

func decode(src []byte) (dst []byte, err error) {
	//return base64.StdEncoding.DecodeString(string(s))
	dst = make([]byte, len(src))
	_, err = base64.StdEncoding.Decode(dst, src)
	return
}
