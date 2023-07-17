package main

import (
	"sync"

	"github.com/PretendoNetwork/friends-secure/grpc"
	"github.com/PretendoNetwork/friends-secure/nex"
)

var wg sync.WaitGroup

func main() {
	wg.Add(3)

	go grpc.StartGRPCServer()
	go nex.StartAuthenticationServer()
	go nex.StartSecureServer()

	wg.Wait()
}
