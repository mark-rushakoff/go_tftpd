package safe_packets

import (
	"bytes"
	"encoding/binary"

	"github.com/mark-rushakoff/go_tftpd/packets"
)

type SafeAck struct {
	packets.Ack
}

func NewSafeAck(blockNumber uint16) *SafeAck {
	return &SafeAck{
		packets.Ack{blockNumber},
	}
}

func (ack *SafeAck) Bytes() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, packets.AckOpcode)
	binary.Write(buf, binary.BigEndian, ack.BlockNumber)
	return buf.Bytes()
}
