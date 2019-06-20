package util

import (
	"encoding/json"
	"io/ioutil"
)

func LoadOptionsFile(optionsFilePath string, options interface{}) error {
	var (
		err     error
		content []byte
	)

	content, err = ioutil.ReadFile(optionsFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, options)
	if err != nil {
		return err
	}

	return nil
}
