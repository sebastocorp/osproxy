package pools

import (
	"crypto/md5"
	"encoding/hex"
	"maps"
	"sync"

	"osproxy/internal/objectStorage"
)

type ActionPoolT struct {
	mu       sync.RWMutex
	requests map[string]ActionPoolRequestT
}

type ActionPoolRequestT struct {
	Object objectStorage.ObjectT
}

func NewActionPool() (p *ActionPoolT) {
	p = &ActionPoolT{
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
	key := hex.EncodeToString(md5.New().Sum([]byte(r.Object.String())))
	p.mu.Lock()
	p.requests[key] = r
	p.mu.Unlock()
}

func (p *ActionPoolT) Remove(key string) {
	p.mu.Lock()
	delete(p.requests, key)
	p.mu.Unlock()
}
