package zerorpc

import (
	"bytes"
	"encoding/gob"
)

type RpcData struct {
	Name string
	Args []interface{}
}

func encode(data RpcData) ([]byte, error) {
	var buf bytes.Buffer
	bufEnc := gob.NewEncoder(&buf)

	//bufEnc.Encode(data)

	if err := bufEnc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(bs []byte) (RpcData, error) {
	buf := bytes.NewBuffer(bs)
	bufDec := gob.NewDecoder(buf)
	var data RpcData
	if err := bufDec.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}
