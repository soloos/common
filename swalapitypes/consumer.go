package swalapitypes

type Consumer interface {
	Consume(msg []byte) error
}
