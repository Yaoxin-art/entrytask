package zerorpc

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestSession_ReadWrite(t *testing.T) {
	addr := "127.0.0.1:6666"
	data := "Hello world."

	wg := sync.WaitGroup{}

	wg.Add(2)

	// write
	go func() {
		defer wg.Done()

		listen, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn, _ := listen.Accept()
		s := Session{conn: conn}

		err = s.Write([]byte(data))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("write success. ")
	}()

	// read
	go func() {
		defer wg.Done()

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}

		s := Session{conn: conn}

		readed, errRead := s.Read()
		if errRead != nil {
			t.Fatal(errRead)
		}
		if string(readed) != data {
			t.Fatal(errRead)
		}
		fmt.Printf("read success: %s \n", string(readed))
	}()

	wg.Wait()
}
