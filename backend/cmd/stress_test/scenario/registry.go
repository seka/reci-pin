package scenario

import (
	"context"
	"net/http"
	"sync"
)

type Scenario func(ctx context.Context, client *http.Client, baseURL string) error

var (
	registryMu sync.RWMutex
	registry   = make(map[string]Scenario)
)

// Register adds a scenario to the registry
func Register(name string, s Scenario) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = s
}

// Get retrieves a scenario from the registry
func Get(name string) (Scenario, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	s, ok := registry[name]
	return s, ok
}
