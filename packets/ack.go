package packets

const AckOpcode uint16 = 4

type Ack struct {
	BlockNumber uint16
}
