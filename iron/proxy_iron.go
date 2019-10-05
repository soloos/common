package iron

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"golang.org/x/xerrors"
)

func (p *Proxy) DispatchWithIronRequest(path string, reqCtx RequestContext, ir *Request) IRespData {
	if !p.IsServiceExists(path) {
		return MakeResp(nil, xerrors.Errorf("%w,path:%s", ErrCmdNotFound, path))
	}

	var reqArgBytes, err = ioutil.ReadAll(ir.R.Body)
	if err != nil {
		return MakeResp(nil, err)
	}

	var service = p.ServiceTable[path]
	var reqArgElems []interface{}

	var parseEasyKvReqArgs = func() {
		var reqArgs = MakeEasyKvReqArgs()

		//merge url params
		reqArgs.MergeIronRequest(ir)

		//merge body params
		if len(reqArgBytes) != 0 {
			var ret = make(map[string]interface{})
			err = json.Unmarshal(reqArgBytes, &ret)
			if err != nil {
				return
			}
			reqArgs.MergeKv(ret)
		}

		reqArgElems = append(reqArgElems, reqArgs)
	}

	var parseNormalReqArgs = func() {
		// parse QueryString
		if service.IsHasUrlKvReqArgs {
			var reqArgs = MakeUrlKvReqArgs()
			reqArgs.MergeIronRequest(ir)
			reqArgElems = append(reqArgElems, reqArgs)
		}

		// parse http body json
		if len(reqArgBytes) == 0 {
			err = ErrCmdParamEmpty
			return
		}

		var reqArgValues []reflect.Value
		var reqArgInterfaces []interface{}
		for i, _ := range service.Params {
			var serviceParam = service.Params[i]
			var reqArgValue = reflect.New(serviceParam)
			reqArgValues = append(reqArgValues, reqArgValue)
			reqArgInterfaces = append(reqArgInterfaces, reqArgValue.Interface())
		}

		if len(reqArgInterfaces) == 1 {
			err = json.Unmarshal(reqArgBytes, &reqArgInterfaces[0])
		} else {
			err = json.Unmarshal(reqArgBytes, &reqArgInterfaces)
		}

		if err != nil {
			return
		}

		for i, _ := range reqArgValues {
			reqArgElems = append(reqArgElems, reqArgValues[i].Elem())
		}
	}

	if err != nil {
		return MakeResp(nil, err)
	}

	// parse EasyKvReqArgs
	if service.IsHasEasyKvReqArgs {
		parseEasyKvReqArgs()
	} else {
		parseNormalReqArgs()
	}

	return p.Dispatch(path, reqCtx, reqArgElems...)
}
