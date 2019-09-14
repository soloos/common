package iron

import (
	"encoding/json"
	"reflect"
)

type LowReqArgs = interface{}

func MarshalLowReqArgs(reqArgs ...LowReqArgs) string {
	var toMarshal []interface{}
	for i, _ := range reqArgs {
		if _, ok := reqArgs[i].(reflect.Value); ok {
			toMarshal = append(toMarshal, reqArgs[i].(reflect.Value).Interface())
		} else {
			toMarshal = append(toMarshal, reqArgs[i])
		}
	}
	var reqArgsBytes, err = json.Marshal(toMarshal)
	if err != nil {
		return err.Error()
	}
	return string(reqArgsBytes)
}
