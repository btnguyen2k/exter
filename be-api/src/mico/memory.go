package mico

import (
	"sync"
	"time"
)

// NewMemoryCache creates a new instance of MemoryCache
func NewMemoryCache(config *CacheConfig) ICache {
	absCache := NewAbstractCache(config)
	return &MemoryCache{AbstractCache: absCache}
}

// MemoryCache is an in-memory, in-process cache
type MemoryCache struct {
	*AbstractCache
	lock sync.RWMutex
	data map[string]*CacheEntry
}

func (c *MemoryCache) ensureInit() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.data == nil {
		c.data = make(map[string]*CacheEntry)
	}
}

// Get implements ICache.Get
func (c *MemoryCache) Get(key string) ([]byte, error) {
	c.ensureInit()
	c.lock.RLock()
	defer c.lock.RUnlock()
	if entry, ok := c.data[key]; ok {
		return append([]byte{}, entry.Data...), nil
	}
	return nil, nil
}

// Set implements ICache.Set
func (c *MemoryCache) Set(key string, value []byte) error {
	if value == nil {
		return c.Remove(key)
	}
	c.ensureInit()
	c.lock.Lock()
	defer c.lock.Unlock()
	now := time.Now()
	c.data[key] = &CacheEntry{
		Key:        key,
		CreatedAt:  now,
		AccessedAt: now,
		ExpiredAt:  now.Add(1 * time.Hour),
		Data:       append([]byte{}, value...),
	}
	return nil
}

// Remove implements ICache.Remove
func (c *MemoryCache) Remove(key string) error {
	c.ensureInit()
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = nil
	delete(c.data, key)
	return nil
}

// Exists implements ICache.Exists
func (c *MemoryCache) Exists(key string) (bool, error) {
	c.ensureInit()
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data[key] != nil, nil
}
