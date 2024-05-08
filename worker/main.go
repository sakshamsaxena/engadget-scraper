package worker

import (
	"io"
	"strings"

	"github.com/sakshamsaxena/engadget-scraper/pool"
)

// Manager will poll for work (read from reader) and send new work to
// the worker pool (collection of scraper instances)
type Manager struct {
	workSource io.Reader
	workerPool *pool.Pool
}

func NewManager(reader io.Reader) *Manager {
	return &Manager{
		workSource: reader,
		workerPool: &pool.Pool{},
	}
}

func (m *Manager) Poll() {
	buf := strings.Builder{}
	for {
		b := make([]byte, 1)
		n, err := m.workSource.Read(b)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 || err == io.EOF {
			break
		}
		if b[0] == '\n' {
			s := buf.String()
			m.workerPool.AddJob(s)
			buf.Reset()
		} else {
			buf.Write(b)
		}
	}
	m.workerPool.Wait()
}
