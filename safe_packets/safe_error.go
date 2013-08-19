package safe_packets

import (
	"bytes"
	"encoding/binary"

	"github.com/mark-rushakoff/go_tftpd/packets"
)

type SafeError struct {
	Code    packets.ErrorCode
	Message string
}

func NewFileNotFoundError() *SafeError {
	return &SafeError{
		packets.FileNotFound,
		"File not found",
	}
}

func (e *SafeError) Bytes() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, packets.ErrorOpcode)
	binary.Write(buf, binary.BigEndian, e.Code)
	buf.WriteString(e.Message)
	buf.WriteByte(0)
	return buf.Bytes()
}
