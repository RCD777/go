package main

import (
	"fmt"
	"time"
)

func main() {

	c := make(chan int)
	c <- 1

	time.AfterFunc(1*time.Second, func() {
		close(c)
	})

	var d chan int

	select {
	case it := <-d:
		fmt.Println("exit")
		fmt.Println(it)

	case c <- 0:

	}

	fmt.Println("ddddd")

	//panic("stack")
}
