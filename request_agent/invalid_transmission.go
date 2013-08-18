package request_agent

type InvalidTransmission struct {
	Packet []byte
	Reason InvalidTransmissionReason
}
