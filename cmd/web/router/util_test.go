package router

import "testing"

func TestString2Bytes(t *testing.T) {
	str := "hello world!"
	bytes := String2Bytes(str)
	t.Logf("string:[%s] to byte[] len:%d", str, len(bytes))

	dist := Bytes2String(bytes)
	t.Logf("bytes to string:[%s", dist)
}
