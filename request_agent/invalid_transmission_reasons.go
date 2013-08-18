package request_agent

type InvalidTransmissionReason int

const (
	PacketTooShort InvalidTransmissionReason = iota
	InvalidOpcode
)
