package solomqapitypes

type Consumer interface {
	Consume(msg []byte) error
}
