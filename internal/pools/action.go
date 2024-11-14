package pools

import (
	"maps"
	"sync"

	"osproxy/internal/objectstorage"
)

type ActionPoolT struct {
	mu       sync.RWMutex
	cap      int
	current  int
	requests map[string]ActionPoolRequestT
}

type ActionPoolRequestT struct {
	Object objectstorage.ObjectT
}

func NewActionPool(cap int) (p *ActionPoolT) {
	p = &ActionPoolT{
		cap:      cap,
		current:  0,
		requests: map[string]ActionPoolRequestT{},
	}

	return p
}

func (p *ActionPoolT) Get() map[string]ActionPoolRequestT {
	result := map[string]ActionPoolRequestT{}

	p.mu.RLock()
	maps.Copy(result, p.requests)
	p.mu.RUnlock()

	return result
}

func (p *ActionPoolT) Add(r ActionPoolRequestT) {
	key := r.Object.StructHash()
	p.mu.Lock()
	if p.current < p.cap {
		p.requests[key] = r
		p.current++
	}
	p.mu.Unlock()
}

func (p *ActionPoolT) Remove(key string) {
	p.mu.Lock()
	delete(p.requests, key)
	p.current--
	p.mu.Unlock()
}
