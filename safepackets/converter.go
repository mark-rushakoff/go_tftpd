package safepackets

import (
	"strings"

	"github.com/mark-rushakoff/go_tftpd/packets"
)

type Converter interface {
	FromAck(ack *packets.Ack) *SafeAck
	FromReadRequest(read *packets.ReadRequest) (*SafeReadRequest, *ConversionError)
}

type converter struct{}

type ConversionError struct {
	code    packets.ErrorCode
	message string
}

func NewConversionError(code packets.ErrorCode, message string) *ConversionError {
	return &ConversionError{
		code:    code,
		message: message,
	}
}

func (e *ConversionError) Error() string {
	return e.message
}

func (e *ConversionError) Code() packets.ErrorCode {
	return e.code
}

func NewConverter() Converter {
	return converter{}
}

func (converter) FromAck(ack *packets.Ack) *SafeAck {
	return &SafeAck{
		packets.Ack{BlockNumber: ack.BlockNumber},
	}
}

func (converter) FromReadRequest(read *packets.ReadRequest) (*SafeReadRequest, *ConversionError) {
	var mode ReadWriteMode
	switch strings.ToLower(read.Mode) {
	case "netascii":
		mode = NetAscii
	case "octet":
		mode = Octet
	default:
		return nil, &ConversionError{code: packets.Undefined, message: "Invalid mode string"}
	}

	return &SafeReadRequest{
		Filename: read.Filename,
		Mode:     mode,
	}, nil
}
