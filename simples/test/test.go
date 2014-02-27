package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 10)

	c <- 1
	c <- 2

	go func() {
		for it := range c {
			fmt.Println(it)
		}

	}()

	fmt.Println("ddddd")
	//close(c)

	time.Sleep(5 * time.Second)

	panic("stack")
}
