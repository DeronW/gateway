package protocol

type Message struct {
	Teleport     int
	DeviceAddr   int
	Op           string
	Status       string
	OriginStatus string
}
