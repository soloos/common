package types

const (
	ServiceIDLen = 128
)

type ServiceID = [ServiceIDLen]byte
type Service func(requestID uint64, conn *Connection)
