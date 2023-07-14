package utils

import (
	"sync"
)

// Mutex A wrapper around sync.Mutex which eliminates the need for constantly calling lock/unlock.
// Very useful to prevent deadlocks, as only top-level functions should be calling RunLocked()
type Mutex struct {
	Mutex *sync.Mutex
}

func NewMutex() *Mutex {
	return &Mutex{Mutex: &sync.Mutex{}}
}

// RunLocked Runs a function in a locked context
func (m *Mutex) RunLocked(f func()) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	f()
}
