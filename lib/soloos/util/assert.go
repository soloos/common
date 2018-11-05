package util

import "log"

func AssertErrIsNil1(ignore interface{}, err error) {
	AssertErrIsNil(err)
}

func AssertNotNil(obj interface{}) {
	if obj == nil {
		log.Panic("obj is nil")
	}
}

func AssertErrIsNil(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

func Ignore(r interface{}) {
}
