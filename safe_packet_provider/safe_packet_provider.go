package safe_packet_provider

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/request_agent"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

type SafePacketProvider struct {
	incomingSafeAck         chan *safety_filter.IncomingSafeAck
	incomingSafeReadRequest chan *safety_filter.IncomingSafeReadRequest
	requestAgent            *request_agent.RequestAgent
}

func NewSafePacketProvider(conn net.PacketConn) *SafePacketProvider {
	ackChan := make(chan *safety_filter.IncomingSafeAck, 3)
	readChan := make(chan *safety_filter.IncomingSafeReadRequest, 3)
	safeRequestHandler := &safeRequestHandler{
		safeAck:         ackChan,
		safeReadRequest: readChan,
	}
	safetyFilter := safety_filter.MakeSafetyFilter(safeRequestHandler)
	requestHandler := &requestHandler{
		safetyFilter: safetyFilter,
	}
	requestAgent := request_agent.NewRequestAgent(conn, requestHandler)

	return &SafePacketProvider{
		incomingSafeAck:         ackChan,
		incomingSafeReadRequest: readChan,
		requestAgent:            requestAgent,
	}
}

func (p *SafePacketProvider) IncomingSafeAck() <-chan *safety_filter.IncomingSafeAck {
	return p.incomingSafeAck
}

func (p *SafePacketProvider) IncomingSafeReadRequest() <-chan *safety_filter.IncomingSafeReadRequest {
	return p.incomingSafeReadRequest
}

func (p *SafePacketProvider) Read() {
	p.requestAgent.Read()
}