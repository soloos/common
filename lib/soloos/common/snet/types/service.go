package types

const (
	ServiceIDLen = 64
)

type ServiceID = [ServiceIDLen]byte
type ServiceRequest struct {
	ReqID        uint64
	ReqBodySize  uint32
	ReqParamSize uint32
	Conn         *Connection
}
type Service func(req ServiceRequest) error
