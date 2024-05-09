package pool

import (
	"net/http"
	"sync"
	"time"

	"github.com/sakshamsaxena/engadget-scraper/scraper"
)

const (
	MaxJobs = 20
)

func New() *Pool {
	return &Pool{
		wg: &sync.WaitGroup{},
		mu: &sync.RWMutex{},
		client: &http.Client{
			Transport: &http.Transport{
				// Setting this to ensure that the client can maintain
				// at least MaxJobs connections so that none of the
				// scraper is has to re-establish connection wherever
				// possible.
				MaxIdleConnsPerHost: MaxJobs,
			},
			Timeout: 2 * time.Second,
		},
		liveJobs: 0,
	}
}

type Pool struct {
	wg       *sync.WaitGroup
	mu       *sync.RWMutex
	client   *http.Client
	liveJobs int
}

func (p *Pool) AddJob(url string) {
	for p.jobs() >= MaxJobs {
		// wait
	}
	p.incr()
	p.wg.Add(1)
	go p.runJob(url)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) decr() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.liveJobs--
}

func (p *Pool) incr() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.liveJobs++
}

func (p *Pool) jobs() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.liveJobs
}

func (p *Pool) runJob(url string) {
	defer p.decr()
	defer p.wg.Done()
	scraper.New(p.client).Scrape(url)
}
