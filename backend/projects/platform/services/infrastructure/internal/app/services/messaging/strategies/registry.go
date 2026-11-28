package strategies

import (
	"fmt"
	"sync"
)

type Registry[T any] struct {
	mu         sync.RWMutex
	strategies map[string]T
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{strategies: make(map[string]T)}
}

func (registry *Registry[T]) Register(key string, value T) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.strategies[key] = value
}

func (registry *Registry[T]) Resolve(key string) (T, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	value, ok := registry.strategies[key]
	if !ok {
		var zero T
		return zero, fmt.Errorf("strategy %q is not registered", key)
	}

	return value, nil
}
