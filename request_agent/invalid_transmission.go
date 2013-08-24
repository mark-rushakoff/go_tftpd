package request_agent

import "net"

type InvalidTransmission struct {
	Packet []byte
	Reason InvalidTransmissionReason
	Addr   net.Addr
}
