package protocol

type Command_login3 struct {
	CommandBase
}

func (c *Command_login3) GetReply() (*PacketSend, bool) {
	return &PacketSend{
		Encrypted:         true,
		WirelessEncrypted: true,
		Op:                4,
		Params:            "",
		Version:           c.Packet.Version,
	}, true
}
