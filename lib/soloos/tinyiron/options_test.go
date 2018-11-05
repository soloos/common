package tinyiron

import "testing"

func TestSanitizeOptions(t *testing.T) {
	var (
		err error
	)

	var server = NewServer()
	err = server.LoadOptionsFile("./options.json")
	AssertErrIsNilForTest(t, err)
}
