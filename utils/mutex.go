package utils

import "sync"

type Mutex struct {
	Mutex *sync.Mutex
}

func NewMutex() *Mutex {
	return &Mutex{Mutex: &sync.Mutex{}}
}

func (m *Mutex) RunLocked(f func()) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	f()
}
