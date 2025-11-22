package provider

import (
	"errors"
	"sync"
)

var (
	mu        sync.RWMutex
	providers = make(map[string]Provider)

	ErrProviderNotFound          = errors.New("could not find the requested provider")
	ErrProviderAlreadyRegistered = errors.New("the provider is already registered")
)

// Registers the requested provider
func Register(p Provider) error {
	mu.Lock()
	defer mu.Unlock()
	name := p.Name()
	if _, exists := providers[name]; exists {
		return ErrProviderAlreadyRegistered
	}
	providers[name] = p
	return nil
}

// Returns the requested provider or an error if not found
func Get(name string) (Provider, error) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := providers[name]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return p, nil
}

// Returns all available providers
func All() []Provider {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]Provider, 0, len(providers))
	for _, p := range providers {
		result = append(result, p)
	}
	return result
}
