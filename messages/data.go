package messages

const DataOpcode uint16 = 3

type Data struct {
	BlockNumber uint16
	Data        []byte
}
