package snetprotocol

import "encoding/gob"

//go:generate msgp

func init() {
	gob.Register(MessageTest0{})
	gob.Register(MessageTest1{})
	gob.Register(MessageTest2{})
}

type MessageTest0 struct {
	Data0 string
}

type MessageTest1 struct {
	Data0 string
	Data1 int
}

type MessageTest2 struct {
	Data0 string
	Data1 int
	Data2 string
}
