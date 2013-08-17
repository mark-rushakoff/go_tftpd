package messages

const ErrorOpcode uint16 = 5

type ErrorCode uint16

const (
	Undefined                    ErrorCode = 0
	FileNotFound                 ErrorCode = 1
	AccessViolation              ErrorCode = 2
	DiskFullOrAllocationExceeded ErrorCode = 3
	IllegalTftpOperation         ErrorCode = 5
	FileAlreadyExists            ErrorCode = 6
	NoSuchUser                   ErrorCode = 7
)

type Error struct {
	Code    ErrorCode
	Message string
}
