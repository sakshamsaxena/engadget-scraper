package main

import (
	"fmt"
	"github.com/sakshamsaxena/engadget-scraper/cache"
	"os"

	"github.com/sakshamsaxena/engadget-scraper/workers"
)

func main() {
	// Get a reader stream of file
	jobSource, openErr := os.OpenFile("endg-urls-test", os.O_RDONLY, 0444)
	if openErr != nil {
		panic(openErr)
	}

	// Initiate cache
	cache.Initialize()

	// Start the workers with this stream
	manager := workers.NewManager(jobSource)

	// Wait for the stream to get over
	<-manager.PollForWork()

	// Dump answer to STDOUT
	fmt.Println(manager.Results())
}
