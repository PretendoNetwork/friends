package main

import (
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	go startGRPCServer()
	go startNEXServer()

	wg.Wait()
}
