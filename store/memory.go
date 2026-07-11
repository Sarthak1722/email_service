package store

import (
	"fmt"
	"sync"

	"github.com/Sarthak1722/email_service/email"
)

type MemoryStore struct {
	mu   sync.RWMutex
	jobs map[string]*email.Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs: make(map[string]*email.Job),
	}
}

func (m *MemoryStore) Save(job *email.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs[job.ID] = job
	return nil
}

func (m *MemoryStore) Get(id string) (*email.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, exists := m.jobs[id]
	if !exists {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}

func (m *MemoryStore) All() []*email.Job {
	m.mu.RLock()
	defer m.mu.RUnlock()
	jobs := make([]*email.Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}