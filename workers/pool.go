package workers

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sakshamsaxena/engadget-scraper/cache"
)

const (
	MaxJobs = 20
)

func newPool() *pool {
	return &pool{
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

type pool struct {
	wg       *sync.WaitGroup
	mu       *sync.RWMutex
	client   *http.Client
	liveJobs int
}

func (p *pool) AddJob(url string) {
	for p.jobs() >= MaxJobs {
		// wait
	}
	p.incr()
	p.wg.Add(1)
	go p.runJob(url)
}

func (p *pool) Wait() {
	p.wg.Wait()
}

func (p *pool) decr() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.liveJobs--
}

func (p *pool) incr() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.liveJobs++
}

func (p *pool) jobs() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.liveJobs
}

func (p *pool) runJob(url string) {
	defer p.decr()
	defer p.wg.Done()
	result := scrapeURL(p.client, url)
	uncleanTokens := strings.Split(result, " ")
	validTokens := tokenize(uncleanTokens)
	err := cache.SetWords(validTokens)
	if err != nil {
		panic(err) // TODO: retry
	}
}
