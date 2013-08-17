package messages

type ReadWriteMode int

const (
	NetAscii ReadWriteMode = iota
	Octet
)

type readWriteRequest struct {
	Filename string
	Mode     ReadWriteMode
}

type ReadRequest struct {
	readWriteRequest
}

type WriteRequest struct {
	readWriteRequest
}
