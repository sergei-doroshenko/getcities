package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup // 1

func ready(w string, sec int) {
	defer wg.Done() // 3
	time.Sleep(time.Duration(sec) * time.Second)
	fmt.Println(w, "is ready!")
}

func main() {

	fmt.Println("I'm waiting")

	wg.Add(1)          // 2
	go ready("Tea", 2) // *

	wg.Add(1)
	go ready("Coffe", 1)

	wg.Wait() // 4

	fmt.Println("main finished")
}
