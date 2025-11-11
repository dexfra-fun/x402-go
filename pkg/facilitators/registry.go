package facilitators

import (
	"fmt"
	"sync"
)

var (
	globalRegistry = &Registry{
		facilitators: make(map[string]*Facilitator),
	}
)

// Registry manages a collection of facilitators.
type Registry struct {
	mu           sync.RWMutex
	facilitators map[string]*Facilitator
}

// Register adds a facilitator to the registry.
// Returns an error if a facilitator with the same ID already exists.
func (r *Registry) Register(f *Facilitator) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.facilitators[f.ID]; exists {
		return fmt.Errorf("facilitator with ID %q already registered", f.ID)
	}

	r.facilitators[f.ID] = f
	return nil
}

// Get retrieves a facilitator by ID string.
// Returns nil if not found.
func (r *Registry) Get(id string) *Facilitator {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.facilitators[id]
}

// List returns all registered facilitators.
// If network is specified, only facilitators supporting that network are returned.
func (r *Registry) List(network Network) []*Facilitator {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Facilitator
	for _, f := range r.facilitators {
		if network == "" {
			result = append(result, f)
		} else if _, hasNetwork := f.Addresses[network]; hasNetwork {
			result = append(result, f)
		}
	}

	return result
}

// Validate checks if all facilitators in the registry have unique IDs.
func (r *Registry) Validate() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	seen := make(map[string]bool)
	for id := range r.facilitators {
		if seen[id] {
			return fmt.Errorf("duplicate facilitator ID: %q", id)
		}
		seen[id] = true
	}

	return nil
}

// Register adds a facilitator to the global registry.
func Register(f *Facilitator) error {
	return globalRegistry.Register(f)
}

// Get retrieves a facilitator from the global registry by string ID.
func Get(id string) *Facilitator {
	return globalRegistry.Get(id)
}

// List returns all facilitators from the global registry.
func List(network Network) []*Facilitator {
	return globalRegistry.List(network)
}

// Validate validates the global registry.
func Validate() error {
	return globalRegistry.Validate()
}
