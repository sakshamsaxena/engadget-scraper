package main

import (
	"github.com/sakshamsaxena/engadget-scraper/worker"
	"os"
)

func main() {
	// Get a reader stream of file
	f, e := os.OpenFile("endg-urls-test", os.O_RDONLY, 0444)
	if e != nil {
		panic(e)
	}
	man := worker.NewManager(f)
	man.Poll()
	// Start the worker poller with this stream
	// Wait for the stream to get over
	// Dump answer from Redis
}
