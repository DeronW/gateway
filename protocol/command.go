package protocol

type Command interface {
	GetOp() string
	GetReply() (*PacketSend, bool)
	GetMessage() (*Message, bool)
	GetEvent() (*Event, bool)
}

type CommandBase struct {
	Packet *PacketReceive
}

func (c *CommandBase) GetOp() string {
	return c.Packet.Op
}

func (c *CommandBase) GetReply() (p *PacketSend, ok bool) {
	return
}

func (c *CommandBase) GetMessage() (m *Message, ok bool) {
	return
}

func (c *CommandBase) GetEvent() (e *Event, ok bool) {
	return
}
