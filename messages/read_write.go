package messages

type ReadWriteMode int

const (
	ReadOpcode  uint16 = 1
	WriteOpcode uint16 = 2
)

const (
	NetAscii ReadWriteMode = iota
	Octet
)

type ReadRequest struct {
	Filename string
	Mode     ReadWriteMode
}

type WriteRequest struct {
	Filename string
	Mode     ReadWriteMode
}
