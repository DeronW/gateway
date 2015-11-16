package teleport

import (
	"gateway/protocol"
	log "github.com/Sirupsen/logrus"
	tcp "github.com/delongw/phantom-tcp"
	"time"
)

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
		log.WithFields(log.Fields{
			"error": err,
		}).Info("can not distatch this request")
	}
	return true
}

func (p *Pool) OnClose(c *tcp.Conn) {
	log.WithFields(log.Fields{
		"remote addr": c.RemoteAddr(),
		"time":        time.Now(),
	}).Info("closed")
}
