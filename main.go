package main

import (
	"os"

	"github.com/sakshamsaxena/engadget-scraper/manager"
)

func main() {
	// Get a reader stream of file
	jobSource, openErr := os.OpenFile("endg-urls-test", os.O_RDONLY, 0444)
	if openErr != nil {
		panic(openErr)
	}

	// Start the workers with this stream
	workManager := manager.New(jobSource)

	// Wait for the stream to get over
	<-workManager.PollForWork()

	// Dump answer from Redis
	// TODO
}
