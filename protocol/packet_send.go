package protocol

type PacketSend struct {
	DeviceAddr        uint32
	Encrypted         bool
	WirelessEncrypted bool
	Op                uint8
	Params            string
	Version           int
}
