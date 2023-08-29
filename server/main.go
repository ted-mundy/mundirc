package main

import (
	"flag"
	"sync"

	"mundy.io/mundirc/pkg/net"
)

var wg sync.WaitGroup

func listen(port int) {
	defer wg.Done()
	net.Listen(port)
}

func main() {
	port := flag.Int("p", 6667, "port to listen on")
	flag.Parse()

	wg.Add(1)
	go listen(*port)
	wg.Wait()
}
