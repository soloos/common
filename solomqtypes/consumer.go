package solomqtypes

type Consumer interface {
	Consume(msg []byte) error
}
