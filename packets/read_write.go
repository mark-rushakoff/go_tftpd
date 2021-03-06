package packets

const (
	ReadOpcode  uint16 = 1
	WriteOpcode uint16 = 2
)

type ReadRequest struct {
	Filename string
	Mode     string
	Options  map[string]string
}

type WriteRequest struct {
	Filename string
	Mode     string
	Options  map[string]string
}
