package zerorpc

import (
	"bytes"
	"encoding/gob"
	"github.com/sirupsen/logrus"
)

type Invocation struct {
	MethodId string
	Args     []interface{}
}

type Result struct {
	StatusCode        int           // status code, 1-success,404-method not exists,0-error
	RequestTimestamp  int64         // get request timestamp
	ResponseTimestamp int64         // out response timestamp
	Err               string        // error if exists
	Args              []interface{} // response body
}

// encode encode zerorpc invocation data
func encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	bufEncoded := gob.NewEncoder(&buf)
	err := bufEncoded.Encode(data)
	if err != nil {
		logrus.Errorf("encode invocation data err:%v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// decodeInvocation decode zerorpc invocation request
func decodeInvocation(data []byte) (Invocation, error) {
	buf := bytes.NewBuffer(data)
	bufDecode := gob.NewDecoder(buf)
	var invocation Invocation
	err := bufDecode.Decode(&invocation)
	if err != nil {
		return Invocation{}, err
	}
	return invocation, nil
}

// decodeResult decode zerorpc invocation response
func decodeResult(data []byte) (Result, error) {
	buf := bytes.NewBuffer(data)
	bufDecode := gob.NewDecoder(buf)
	var result Result
	err := bufDecode.Decode(&result)
	if err != nil {
		return Result{}, err
	}
	return result, nil
}
