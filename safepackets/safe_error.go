package safepackets

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
		Code:    packets.FileNotFound,
		Message: "File not found",
	}
}

func NewAccessViolationError(message string) *SafeError {
	return &SafeError{
		Code:    packets.AccessViolation,
		Message: message,
	}
}

func (e *SafeError) Equals(other *SafeError) bool {
	return e.Code == other.Code && e.Message == other.Message
}

func (e *SafeError) Bytes() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, packets.ErrorOpcode)
	binary.Write(buf, binary.BigEndian, e.Code)
	buf.WriteString(e.Message)
	buf.WriteByte(0)
	return buf.Bytes()
}
