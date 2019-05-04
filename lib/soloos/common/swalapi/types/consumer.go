package types

type Consumer interface {
	Consume(msg []byte) error
}
