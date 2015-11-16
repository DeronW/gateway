package protocol

type Command interface {
	GetOp() string
	GetSendPacket() (*PacketSend, bool)
	GetPublishMessage() (*Message, bool)
}
