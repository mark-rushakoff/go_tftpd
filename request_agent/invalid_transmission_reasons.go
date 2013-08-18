package request_agent

import "fmt"

type InvalidTransmissionReason int

const (
	PacketTooShort InvalidTransmissionReason = iota
	InvalidOpcode
)

func (reason InvalidTransmissionReason) String() string {
	switch reason {
	case PacketTooShort:
		return "Packet too short"
	case InvalidOpcode:
		return "Invalid opcode"
	default:
		panic(fmt.Sprintf("No string exists for reason code %d", reason))
	}
}
