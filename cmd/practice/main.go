package main

import (
	"container/list"
	"fmt"
)

func main() {
	users := list.New()
	users.PushBack(23)
	users.PushBack("good")
	users.PushBack(true)

	for e := users.Front(); e != nil; e = e.Next() {
		fmt.Printf("e:%v \n", e.Value)
	}
}
