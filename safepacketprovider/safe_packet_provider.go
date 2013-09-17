package safepacketprovider

import (
	"net"

	"github.com/mark-rushakoff/go_tftpd/requestagent"
	"github.com/mark-rushakoff/go_tftpd/safepackets"
	"github.com/mark-rushakoff/go_tftpd/safetyfilter"
)

type SafePacketProvider struct {
	incomingSafeAck         chan *safetyfilter.IncomingSafeAck
	incomingSafeReadRequest chan *safetyfilter.IncomingSafeReadRequest
	incomingInvalidMessage  chan *safetyfilter.IncomingInvalidMessage
	requestAgent            *requestagent.RequestAgent
}

func NewSafePacketProvider(conn net.PacketConn) *SafePacketProvider {
	ackChan := make(chan *safetyfilter.IncomingSafeAck, 3)
	readChan := make(chan *safetyfilter.IncomingSafeReadRequest, 3)
	invalidChan := make(chan *safetyfilter.IncomingInvalidMessage, 3)
	safeRequestHandler := &safeRequestHandler{
		safeAck:            ackChan,
		safeReadRequest:    readChan,
		safeInvalidMessage: invalidChan,
	}
	safetyFilter := safetyfilter.MakeSafetyFilter(safepackets.NewConverter(), safeRequestHandler)
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

func (p *SafePacketProvider) IncomingSafeAck() <-chan *safetyfilter.IncomingSafeAck {
	return p.incomingSafeAck
}

func (p *SafePacketProvider) IncomingSafeReadRequest() <-chan *safetyfilter.IncomingSafeReadRequest {
	return p.incomingSafeReadRequest
}

func (p *SafePacketProvider) IncomingInvalidMessage() <-chan *safetyfilter.IncomingInvalidMessage {
	return p.incomingInvalidMessage
}

func (p *SafePacketProvider) Read() {
	p.requestAgent.Read()
}
