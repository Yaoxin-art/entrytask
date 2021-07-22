package router

import "testing"

func TestMd5String(t *testing.T) {
	word := "hello word!"
	t.Logf("md5:%s result: %s", word, Md5String(word))
	word = "Good night"
	t.Logf("md5:%s result: %s", word, Md5String(word))
}
