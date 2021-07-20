package zerorpc

import (
	"fmt"
	"reflect"
	"testing"
)

const serverAddr = ":9999"

func SayHelloAlia(word string) string {
	res := "hello, " + word
	return res
}

func TestClient_CallRPC(t *testing.T) {
	var sum func(a, b int) int
	fmt.Printf("sum type: %T", sum)
	server := NewServer(serverAddr)
	server.Register("sayHello", SayHelloAlia)
	go server.Run()

	var sayHello func(word string) string
	inter := &sayHello
	fn := reflect.ValueOf(inter).Elem()
	fmt.Println(fn.Kind())

	client := NewClient()
	client.Config("sayHello", serverAddr, &sayHello)

	res := sayHello("word")
	t.Logf("test res:%v", res)
}
