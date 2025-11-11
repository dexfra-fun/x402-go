package x402

import (
	"sync"
	"time"
)

// FeePayerCache caches fee payer addresses to reduce facilitator calls
type FeePayerCache struct {
	mu       sync.RWMutex
	data     map[string]cachedFeePayer
	cacheTTL time.Duration
}

type cachedFeePayer struct {
	feePayer  string
	timestamp time.Time
}

// NewFeePayerCache creates a new fee payer cache
func NewFeePayerCache(ttl time.Duration) *FeePayerCache {
	return &FeePayerCache{
		data:     make(map[string]cachedFeePayer),
		cacheTTL: ttl,
	}
}

// Get retrieves a cached fee payer if it exists and hasn't expired
func (c *FeePayerCache) Get(network string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cached, ok := c.data[network]; ok {
		if time.Since(cached.timestamp) < c.cacheTTL {
			return cached.feePayer, true
		}
		// Expired, will be removed on next cleanup
	}
	return "", false
}

// Set stores a fee payer in the cache
func (c *FeePayerCache) Set(network, feePayer string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[network] = cachedFeePayer{
		feePayer:  feePayer,
		timestamp: time.Now(),
	}
}

// Clear removes all entries from the cache
func (c *FeePayerCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]cachedFeePayer)
}

// CleanupExpired removes expired entries from the cache
func (c *FeePayerCache) CleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for network, cached := range c.data {
		if now.Sub(cached.timestamp) >= c.cacheTTL {
			delete(c.data, network)
		}
	}
}

// StartCleanupRoutine starts a background goroutine to cleanup expired entries
func (c *FeePayerCache) StartCleanupRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.CleanupExpired()
		}
	}()
}
