package teleport

import (
	tcp "github.com/delongw/phantom_tcp"
)

type Pool struct {
	conns map[int]*tcp.Conn
}

func (p *Pool) OnConnect(c *tcp.Conn) bool {
	return true
}

func (p *Pool) OnMessage(c *tcp.Conn, m []byte) bool {
	return true
}

func (p *Pool) OnClose(c *tcp.Conn) {
	if c.Id != nil {
		//conns[c.id]
	}
}
