package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 10)

	time.AfterFunc(1*time.Second, func() {
		c <- 1
		c <- 2
	})

	time.AfterFunc(2*time.Second, func() {
		close(c)

	})

	for it := range c {
		fmt.Println(it)
	}

	fmt.Println("ddddd")

	panic("stack")
}
