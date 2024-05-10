package processor

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/sakshamsaxena/engadget-scraper/cache"
)

// Processor will poll for work (read from reader) and send new work to
// the processor pool to run scrapers in multiple goroutines.
type Processor struct {
	workSource io.Reader
	workerPool *pool
}

func New(source io.Reader) *Processor {
	return &Processor{
		workSource: source,
		workerPool: newPool(),
	}
}

func (p *Processor) PollForWork() chan bool {
	done := make(chan bool)
	go p.poll(done)
	return done
}

func (p *Processor) Results() string {
	topTenWords := cache.GetTopNWords(10)
	type Result struct {
		Word string `json:"word"`
		Freq int    `json:"occurrences"`
	}
	type Results struct {
		Results []Result `json:"results"`
	}
	results := Results{Results: make([]Result, 10)}
	for rank, member := range topTenWords {
		results.Results[rank] = Result{
			Freq: (member["freq"]).(int),
			Word: member["word"].(string),
		}
	}
	jsonResults, _ := json.MarshalIndent(results, "", "\t")
	return string(jsonResults)
}

func (p *Processor) poll(done chan bool) {
	buffer := strings.Builder{}
	for {
		singleByte := make([]byte, 1)
		bytesRead, readErr := p.workSource.Read(singleByte)
		if readErr != nil && readErr != io.EOF {
			panic(readErr)
		}
		if bytesRead == 0 || readErr == io.EOF {
			break
		}
		if singleByte[0] == '\n' {
			// Send the parsed URL to workerPool. It is important
			// to stay blocked here to not spawn too many goroutines
			// that will have to wait. The current goroutine itself
			// will act as the flow manager of tasks to the workerPool
			// by simply honouring workerPool's state.
			url := buffer.String()
			p.workerPool.AddJob(url)
			fmt.Printf("Queued URL %s for scraping\n", url)
			buffer.Reset()
		} else {
			buffer.Write(singleByte)
		}
	}
	if buffer.Len() > 0 {
		url := buffer.String()
		p.workerPool.AddJob(url)
		fmt.Printf("Queued URL %s for scraping\n", url)
		buffer.Reset()
	}
	p.workerPool.Wait()
	done <- true
}
