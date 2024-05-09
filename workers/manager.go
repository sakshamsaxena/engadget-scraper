package workers

import (
	"io"
	"strings"

	"github.com/sakshamsaxena/engadget-scraper/cache"
)

// Manager will poll for work (read from reader) and send new work to
// the worker pool to run scrapers in multiple goroutines.
type Manager struct {
	workSource io.Reader
	workerPool *pool
}

func NewManager(source io.Reader) *Manager {
	return &Manager{
		workSource: source,
		workerPool: newPool(),
	}
}

func (m *Manager) PollForWork() chan bool {
	done := make(chan bool)
	go m.poll(done)
	return done
}

func (m *Manager) Results() string {
	foo := cache.GetTopNWords()
	return strings.Join(foo, "")
}

func (m *Manager) poll(done chan bool) {
	buffer := strings.Builder{}
	for {
		singleByte := make([]byte, 1)
		bytesRead, readErr := m.workSource.Read(singleByte)
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
			m.workerPool.AddJob(buffer.String())
			buffer.Reset()
		} else {
			buffer.Write(singleByte)
		}
	}
	if buffer.Len() > 0 {
		m.workerPool.AddJob(buffer.String())
		buffer.Reset()
	}
	m.workerPool.Wait()
	done <- true
}
