package mico

import (
	"time"
)

const (
	// Version number of mico
	Version = "0.1.0"
)

// CacheEntry represents an entry in cache
type CacheEntry struct {
	Key        string    `json:"key"`
	CreatedAt  time.Time `json:"t_created"`
	AccessedAt time.Time `json:"t_accessed"`
	ExpiredAt  time.Time `json:"t_expired"`
	Data       []byte    `json:"data"`
}

// CacheConfig specifies how cache should behave
type CacheConfig struct {
}

// ICache defines
type ICache interface {
	// Get fetches an entry from cache
	Get(key string) ([]byte, error)

	// Set stores an entry to cache
	Set(key string, value []byte) error

	// Remove deletes an entry from cache
	Remove(key string) error

	// Exists checks if an entry exists in cache
	Exists(key string) (bool, error)
}

// AbstractCache is abstract implementation of cache.
type AbstractCache struct {
	config *CacheConfig
}

// NewAbstractCache creates a new instance of AbstractCache
func NewAbstractCache(config *CacheConfig) *AbstractCache {
	var cloneConf CacheConfig
	if config != nil {
		cloneConf = *config
	}
	return &AbstractCache{&cloneConf}
}

// GetConfig returns cache's configuration settings
func (c *AbstractCache) GetConfig() CacheConfig {
	return *c.config
}
