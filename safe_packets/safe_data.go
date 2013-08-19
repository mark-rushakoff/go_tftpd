package safe_packets

import (
	"bytes"
	"encoding/binary"

	"github.com/mark-rushakoff/go_tftpd/packets"
)

type SafeData struct {
	packets.Data
}

func NewSafeData(blockNumber uint16, data []byte) *SafeData {
	return &SafeData{
		packets.Data{blockNumber, data},
	}
}

func (data *SafeData) Bytes() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, packets.DataOpcode)
	binary.Write(buf, binary.BigEndian, data.BlockNumber)
	buf.Write(data.Data.Data)
	return buf.Bytes()
}
