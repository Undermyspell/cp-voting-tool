package models

import "sync"

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (safeCounter *SafeCounter) Increment() {
	safeCounter.mu.Lock()
	defer safeCounter.mu.Unlock()
	safeCounter.value++
}

func (safeCounter *SafeCounter) Decrement() {
	safeCounter.mu.Lock()
	defer safeCounter.mu.Unlock()
	safeCounter.value--
}

func (safeCounter *SafeCounter) Value() int {
	safeCounter.mu.Lock()
	defer safeCounter.mu.Unlock()
	return safeCounter.value
}
