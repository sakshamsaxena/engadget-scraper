package main

import (
	"fmt"
	"os"

	"github.com/sakshamsaxena/engadget-scraper/cache"
	"github.com/sakshamsaxena/engadget-scraper/config"
	"github.com/sakshamsaxena/engadget-scraper/processor"
)

func main() {
	// Get a reader stream of file
	jobSource, openErr := os.OpenFile(config.Get("static.linkURLs").(string), os.O_RDONLY, 0444)
	if openErr != nil {
		panic(openErr)
	}

	// Initiate cache
	cache.Initialize()

	// Start the processor with this stream
	jobManager := processor.New(jobSource)

	// Wait for the stream to get over
	<-jobManager.PollForWork()

	// Dump answer to STDOUT
	fmt.Println(jobManager.Results())
}
