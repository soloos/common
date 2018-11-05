package snet

type FixedMessageHeader struct {
	Address [128]byte
	Service [256]byte
}
