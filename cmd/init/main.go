package main

import "fmt"

func main() {
	practiceChan()
	practiceSelectChan()
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
