package types

const (
	ServiceIDLen = 64
)

type ServiceID = [ServiceIDLen]byte
type Service func(requestID uint64, requestContentLen uint32, conn *Connection) error
