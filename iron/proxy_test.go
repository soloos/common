package iron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testProxyUrlP(port int, path string) string {
	urlPath := fmt.Sprintf("http://localhost:%v/Argon%v", port, path)
	log.Println(urlPath)
	return urlPath
}

func prepareServer(port int) {
	var testProxyListenStr = fmt.Sprintf("0.0.0.0:%v", port)
	var testProxyServeStr = fmt.Sprintf("http://localhost:%v", port)

	var proxy Proxy
	AssertErrIsNil(proxy.Init())
	proxy.RegisterService("/Test", ProxyServiceBase)
	proxy.RegisterService("/TestMultiArg", ProxyServiceTestMultiArg)
	proxy.RegisterService("/TestLowReqArgs", ProxyServiceTestLowReqArgs)
	proxy.RegisterService("/TestUrlKvReqArgs", ProxyServiceTestUrlKvReqArgs)
	proxy.RegisterService("/TestEasyKvReqArgs", ProxyServiceTestEasyKvReqArgs)

	var webOptions Options
	webOptions.ListenStr = testProxyListenStr
	webOptions.ServeStr = testProxyServeStr
	proxy.InitStandAloneWebServer("/Argon", webOptions)

	go func() {
		AssertErrIsNil(proxy.StandAloneWebServerServe())
	}()
	time.Sleep(time.Millisecond * 200)
}

type ProxyServiceTestReq struct {
	A int
	B int
	C string
}

func ProxyServiceTestLowReqArgs(reqCtx *RequestContext, req0 UrlKvReqArgs) (resp Response) {
	AssertTrue(req0.MustFormFloat64("a", 0.0) == -123.45)
	AssertTrue(req0.MustFormString("b", "error") == "test")
	return Response{fmt.Sprintf("%v", int(req0.MustFormFloat64("a", 0.0))), nil}
}

func TestProxyLowReqArgs(t *testing.T) {
	var serverPort = 17200
	prepareServer(serverPort)
	var (
		resp      *http.Response
		respBytes []byte
		err       error
	)
	assert.NoError(t, err)
	resp, err = http.Post(testProxyUrlP(serverPort, "/TestLowReqArgs?a=-123.45&b=test"),
		"application/json", nil)
	assert.NoError(t, err)
	respBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":"-123","code":0,"message":""}`, string(respBytes))
}

func ProxyServiceTestMultiArg(reqCtx *RequestContext, req0 ProxyServiceTestReq, req1 int) (resp Response) {
	req0.A = req0.A * req1
	req0.B = req0.B * req1
	req0.C = "response"
	return Response{req0, nil}
}

func TestProxyMultiArg(t *testing.T) {
	var serverPort = 17201
	prepareServer(serverPort)
	var reqArgs [2]interface{}
	reqArgs[0] = map[string]interface{}{
		"A": 123,
		"B": 321,
		"C": "helloworld",
	}
	reqArgs[1] = 10
	var (
		reqBytes  []byte
		resp      *http.Response
		respBytes []byte
		err       error
	)
	reqBytes, err = json.Marshal(reqArgs)
	assert.NoError(t, err)
	resp, err = http.Post(testProxyUrlP(serverPort, "/TestMultiArg"),
		"application/json", bytes.NewBuffer(reqBytes))
	assert.NoError(t, err)
	respBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":{"A":1230,"B":3210,"C":"response"},"code":0,"message":""}`, string(respBytes))
}

func ProxyServiceBase(reqCtx *RequestContext, req ProxyServiceTestReq) (resp Response) {
	req.A = req.A * 100
	req.B = req.B * 100
	req.C = "response"
	return Response{req, nil}
}

func TestProxyBase(t *testing.T) {
	var serverPort = 17202
	prepareServer(serverPort)
	var reqArgs = map[string]interface{}{
		"A": 123,
		"B": 321,
		"C": "helloworld",
	}
	var (
		reqBytes  []byte
		resp      *http.Response
		respBytes []byte
		err       error
	)
	reqBytes, err = json.Marshal(reqArgs)
	assert.NoError(t, err)
	resp, err = http.Post(testProxyUrlP(serverPort, "/Test"),
		"application/json", bytes.NewBuffer(reqBytes))
	assert.NoError(t, err)
	respBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":{"A":12300,"B":32100,"C":"response"},"code":0,"message":""}`, string(respBytes))
}

func ProxyServiceTestUrlKvReqArgs(reqCtx *RequestContext, req0 UrlKvReqArgs,
	req1 ProxyServiceTestReq,
	req2 int,
) (resp Response) {
	AssertTrue(req0.MustFormFloat64("a", 0.0) == -123.45)
	AssertTrue(req0.MustFormString("b", "error") == "test")
	return Response{fmt.Sprintf("%v%v", int(req0.MustFormFloat64("a", 0.0))+req1.A, req1.C), nil}
}

func TestProxyUrlKvReqArgs(t *testing.T) {
	var serverPort = 17203
	prepareServer(serverPort)
	var reqArgs [2]interface{}
	reqArgs[0] = map[string]interface{}{
		"A": 10,
		"B": 10,
		"C": "helloworld",
	}
	reqArgs[1] = 10
	var (
		reqBytes  []byte
		resp      *http.Response
		respBytes []byte
		err       error
	)
	reqBytes, err = json.Marshal(reqArgs)
	assert.NoError(t, err)
	resp, err = http.Post(testProxyUrlP(serverPort, "/TestUrlKvReqArgs?a=-123.45&b=test"),
		"application/json", bytes.NewBuffer(reqBytes))
	assert.NoError(t, err)
	respBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":"-113helloworld","code":0,"message":""}`, string(respBytes))
}

func ProxyServiceTestEasyKvReqArgs(reqCtx *RequestContext, req0 EasyKvReqArgs) (resp Response) {
	var a = req0.MustFormFloat64("A", 0.0)
	var b = req0.MustFormFloat64("B", 0.0)
	var c = req0.MustFormString("C", "helloworld")
	AssertTrue(a == 10.0)
	AssertTrue(b == 10.0)
	AssertTrue(c == "test")
	return Response{fmt.Sprintf("%v%v%v", a, b, c), nil}
}

func TestProxyEasyKvReqArgs(t *testing.T) {
	var serverPort = 17204
	prepareServer(serverPort)
	var reqArgs = map[string]interface{}{
		"A": 10.0,
		"C": "test",
	}
	var (
		reqBytes  []byte
		resp      *http.Response
		respBytes []byte
		err       error
	)
	reqBytes, err = json.Marshal(reqArgs)
	assert.NoError(t, err)
	resp, err = http.Post(testProxyUrlP(serverPort, "/TestEasyKvReqArgs?A=30&B=10.0"),
		"application/json", bytes.NewBuffer(reqBytes))
	assert.NoError(t, err)
	respBytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"data":"1010test","code":0,"message":""}`, string(respBytes))
}
