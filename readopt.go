package caches

import (
	"sync"
	"time"
)

// readOptimizedCache sync.map based cache, works best for no-overwrites, only growing caches.
type readOptimizedCache[K comparable] struct {
	storage     sync.Map
	timeout     time.Duration
	cleanupFreq time.Duration
}

func newReadOptimized[K comparable](o cacheOpts) *readOptimizedCache[K] {

	c := &readOptimizedCache[K]{
		timeout:     o.timeout,
		cleanupFreq: o.cleanupFreq,
	}

	if o.cleanupEnabled {
		go c.schedCleanup()
	}

	return c
}

func (c *readOptimizedCache[K]) Put(key K, value any) {
	c.storage.Store(key, item{val: value, until: time.Now().UTC().Add(c.timeout)})
}

func (c *readOptimizedCache[K]) Get(key K) (any, error) {
	entry, found := c.storage.Load(key)
	if !found {
		return nil, ErrNotFound
	}

	item := entry.(item)
	if item.until.Before(time.Now().UTC()) {
		return nil, ErrNotFound
	}

	return item.val, nil
}

func (c *readOptimizedCache[K]) Delete(key K) {
	c.storage.Delete(key)
}

func (c *readOptimizedCache[K]) Drop() {
	c.storage.Clear()
}

func (c *readOptimizedCache[K]) schedCleanup() {
	for {
		<-time.After(c.cleanupFreq)
		c.cleanup()
	}
}

func (c *readOptimizedCache[K]) cleanup() {
	now := time.Now().UTC()

	c.storage.Range(func(key, value any) bool {
		entry := value.(item)

		if now.After(entry.until) {
			c.storage.Delete(key)
		}

		return true
	})

}
