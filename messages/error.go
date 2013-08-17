package messages

type ErrorCode uint16
type ErrorMessage []byte

type Error struct {
	Code    ErrorCode
	Message ErrorMessage
}
