package caches

import (
	"sync"
	"time"
)


// Cache usual map-based cache, well performing in most cases
type Cache[K comparable] struct {
	timeout     time.Duration
	cleanupFreq time.Duration

	storage map[K]*item
	sync.RWMutex
}

func newRW[K comparable](o cacheOpts) *Cache[K] {

	c := &Cache[K]{
		timeout:     o.timeout,
		cleanupFreq: o.cleanupFreq,
		storage:     make(map[K]*item),
	}

	if o.cleanupEnabled {
		go c.schedCleanup()
	}
	return c
}

func (c *Cache[K]) Put(key K, value any) {
	expiryTime := time.Now().UTC().Add(c.timeout)

	c.Lock()
	defer c.Unlock()

	c.storage[key] = &item{
		val:   value,
		until: expiryTime,
	}
}

func (c *Cache[K]) PutWithTimeout(key K, value any, timeout time.Duration) {
	expiryTime := time.Now().UTC().Add(timeout)

	c.Lock()
	defer c.Unlock()

	c.storage[key] = &item{
		val:   value,
		until: expiryTime,
	}
}

func (c *Cache[K]) Get(key K) (any, error) {
	now := time.Now().UTC()

	c.RLock()
	defer c.RUnlock()

	entry, found := c.storage[key]
	if !found {
		return nil, ErrNotFound
	}

	if entry.until.Before(now) {
		return nil, ErrNotFound
	}

	return entry.val, nil
}

func (c *Cache[K]) schedCleanup() {
	for {
		<-time.After(c.cleanupFreq)

		c.cleanup(c.findExpired())
	}
}

func (c *Cache[K]) cleanup(expired []K) {
	c.Lock()
	defer c.Unlock()

	for _, key := range expired {
		delete(c.storage, key)
	}
}

func (c *Cache[K]) Delete(key K) {
	c.Lock()
	defer c.Unlock()

	delete(c.storage, key)
}

func (c *Cache[K]) Drop() {
	c.Lock()
	defer c.Unlock()

	for k, _ := range c.storage {
		delete(c.storage, k)
	}
}

func (c *Cache[K]) findExpired() []K {
	var entries []K

	now := time.Now().UTC()

	c.RLock()
	defer c.RUnlock()

	for k, v := range c.storage {
		if v.until.Before(now) {
			entries = append(entries, k)
		}
	}

	return entries
}
