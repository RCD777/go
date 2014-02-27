package main

import (
	"fmt"
	"time"
)

type Ball struct {
	hits int
}

func Player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball

		//panic("show me the stacks")
	}
}

func main() {
	table := make(chan *Ball)

	go Player("ping", table)
	go Player("pong", table)

	table <- &Ball{0}
	time.Sleep(1 * time.Second)
	<-table //game over, grab the ball

}
