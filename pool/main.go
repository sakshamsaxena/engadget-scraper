package pool

import (
	"fmt"
	"sync"
)

const (
	MaxJobs = 20
)

type Pool struct {
	wg     sync.WaitGroup
	active int
}

func (p *Pool) AddJob(url string) {
	for p.active >= MaxJobs {
		// wait
	}
	p.incr()
	p.wg.Add(1)
	go func(wg *sync.WaitGroup, url string) {
		defer wg.Done()
		defer p.decr()
		// actual scraper i guess ?
		fmt.Println("Starting to scrape ...")
		fmt.Println(url)
		fmt.Println("URL Scraped")
	}(&p.wg, url)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) decr() {
	p.active--
}

func (p *Pool) incr() {
	p.active++
}
