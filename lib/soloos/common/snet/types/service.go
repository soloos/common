package types

const (
	ServiceIDLen = 64
)

type ServiceID = [ServiceIDLen]byte
type Service func(req *NetQuery) error
