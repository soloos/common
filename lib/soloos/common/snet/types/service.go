package types

const (
	ServiceIDLen = 64
)

type ServiceID = [ServiceIDLen]byte
type Service func(reqID uint64, reqBodySize, reqParamSize uint32, conn *Connection) error
