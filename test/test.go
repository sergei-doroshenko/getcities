package main

import (
	"fmt"
)

var done = make(chan bool)
var msgs = make(chan int)

func produce(name string) {
	for i := 0; i < 5; i++ {
		fmt.Println("producer "+name+" made: ", i)
		msgs <- i
	}
	// done <- true
}

func consume(count int) {
	for count > 0 {
		msg := <-msgs
		fmt.Println("consumer get: ", msg)
		count--
	}

	done <- true
}

func main() {
	go produce("1st")
	go produce("2nd")

	go consume(7)
	<-done
}
