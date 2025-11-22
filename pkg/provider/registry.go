package provider

import "sync"

var (
	mu        sync.RWMutex
	providers = make(map[string]Provider)
)

func Register(p Provider) {
	mu.Lock()
	defer mu.Unlock()
	name := p.Name()
	if _, exists := providers[name]; exists {
		panic("provider already registered: " + name)
	}
	providers[name] = p
}

func Get(name string) (Provider, bool) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := providers[name]
	return p, ok
}

func All() []Provider {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]Provider, 0, len(providers))
	for _, p := range providers {
		result = append(result, p)
	}
	return result
}
