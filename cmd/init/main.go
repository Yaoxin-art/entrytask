package main

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"strings"
)

func main() {
	practiceChan()
	practiceSelectChan()

	practiceReflect()

	mysqlPassword("123456")
}

// mysqlPassword mysql password encrypt
// '*' + UPPER(SHA1(UNHEX(SHA1(word))))
func mysqlPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	bs := h.Sum(nil)
	fmt.Println(string(bs))
	bss := string(bs)
	h = sha1.New()
	h.Write([]byte(bss))
	cs := h.Sum(nil)
	css := string(cs)
	fmt.Printf("css upper:%s", strings.ToUpper(css))
	return string(bs)
}

func practiceReflect() {
	var invertInts func([]int) []int
	Bind(&invertInts, InvertSlice)
	fmt.Println(invertInts([]int{1, 2, 3, 4, 2, 3, 5}))
}

func InvertSlice(args []reflect.Value) (result []reflect.Value) {
	inSlice, n := args[0], args[0].Len()
	outSlice := reflect.MakeSlice(inSlice.Type(), 0, n)
	for i := n - 1; i >= 0; i-- {
		element := inSlice.Index(i)
		outSlice = reflect.Append(outSlice, element)
	}
	return []reflect.Value{outSlice}
}

func Bind(p interface{}, f func([]reflect.Value) []reflect.Value) {

	invert := reflect.ValueOf(p).Elem()

	//Use of MakeFunc() method
	invert.Set(reflect.MakeFunc(invert.Type(), f))
}

func practiceChan() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := make(chan int)
	go parallelSum(arr[:len(arr)/2], c)
	go parallelSum(arr[len(arr)/2:], c)
	sum1, sum2 := <-c, <-c
	fmt.Printf("sum: %d \n", sum1+sum2)
}

func parallelSum(arr []int, c chan int) {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	c <- sum
}

func practiceSelectChan() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
