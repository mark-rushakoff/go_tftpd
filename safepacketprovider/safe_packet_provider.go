package safepacketprovider

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/requestagent"
	"github.com/mark-rushakoff/go_tftpd/safety_filter"
)

type SafePacketProvider struct {
	incomingSafeAck         chan *safety_filter.IncomingSafeAck
	incomingSafeReadRequest chan *safety_filter.IncomingSafeReadRequest
	incomingInvalidMessage  chan *safety_filter.IncomingInvalidMessage
	requestAgent            *requestagent.RequestAgent
}

func NewSafePacketProvider(conn net.PacketConn) *SafePacketProvider {
	ackChan := make(chan *safety_filter.IncomingSafeAck, 3)
	readChan := make(chan *safety_filter.IncomingSafeReadRequest, 3)
	invalidChan := make(chan *safety_filter.IncomingInvalidMessage, 3)
	safeRequestHandler := &safeRequestHandler{
		safeAck:            ackChan,
		safeReadRequest:    readChan,
		safeInvalidMessage: invalidChan,
	}
	safetyFilter := safety_filter.MakeSafetyFilter(safeRequestHandler)
	requestHandler := &requestHandler{
		safetyFilter: safetyFilter,
	}
	requestAgent := requestagent.NewRequestAgent(conn, requestHandler)

	return &SafePacketProvider{
		incomingSafeAck:         ackChan,
		incomingSafeReadRequest: readChan,
		incomingInvalidMessage:  invalidChan,
		requestAgent:            requestAgent,
	}
}

func (p *SafePacketProvider) IncomingSafeAck() <-chan *safety_filter.IncomingSafeAck {
	return p.incomingSafeAck
}

func (p *SafePacketProvider) IncomingSafeReadRequest() <-chan *safety_filter.IncomingSafeReadRequest {
	return p.incomingSafeReadRequest
}

func (p *SafePacketProvider) IncomingInvalidMessage() <-chan *safety_filter.IncomingInvalidMessage {
	return p.incomingInvalidMessage
}

func (p *SafePacketProvider) Read() {
	p.requestAgent.Read()
}
