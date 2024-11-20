package caches

import (
	"sync"
)

// lRUCache has good performance for keeping only frequently used values
type lRUCache[K comparable] struct {
	storage map[K]*lruItem[K]
	cap     int
	head    *lruItem[K]
	tail    *lruItem[K]

	sync.RWMutex
}

func newLRU[K comparable](o cacheOpts) *lRUCache[K] {
	return &lRUCache[K]{
		storage: make(map[K]*lruItem[K], o.lruSize),
		cap:     o.lruSize,
	}
}

func (c *lRUCache[K]) moveToFront(entry *lruItem[K]) {
	if c.head == entry {
		return
	}

	if entry.prev != nil {
		entry.prev.next = entry.next
	}

	if entry.next != nil {
		entry.next.prev = entry.prev
	}

	if c.tail == entry {
		c.tail = entry.prev
		if c.tail != nil {
			c.tail.next = nil
		}
	}

	entry.next = c.head
	entry.prev = nil

	if c.head != nil {
		c.head.prev = entry
	}

	c.head = entry

	if c.tail == nil {
		c.tail = entry
	}
}

func (c *lRUCache[K]) Get(key K) (any, error) {
	c.Lock()
	defer c.Unlock()

	entry, found := c.storage[key]
	if !found {
		return nil, ErrKeyExists
	}

	c.moveToFront(entry)

	return entry.value, nil
}

func (c *lRUCache[K]) Put(key K, value any) {
	if entry, found := c.storage[key]; found {
		entry.value = value

		c.moveToFront(entry)

		return
	}

	entry := &lruItem[K]{key: key, value: value}
	c.storage[key] = entry
	c.moveToFront(entry)

	if len(c.storage) > c.cap {
		c.removeTail()
	}
}

func (c *lRUCache[K]) removeTail() {
	if c.tail == nil {
		return
	}

	delete(c.storage, c.tail.key)
	if c.tail.prev != nil {
		c.tail.prev.next = nil
	}

	c.tail = c.tail.prev
	if c.tail == nil {
		c.head = nil
	}
}

func (c *lRUCache[K]) Delete(key K) {
	c.Lock()
	defer c.Unlock()

	entry, found := c.storage[key]
	if !found {
		return
	}

	delete(c.storage, key)

	c.removeFromListWitoutLock(entry)

}

func (c *lRUCache[K]) removeFromListWitoutLock(item *lruItem[K]) {
	if item == nil {
		return
	}
	if item.prev != nil {
		item.prev.next = item.next
	} else {
		c.head = item.next
	}

	if item.next != nil {
		item.next.prev = item.prev
	} else {
		c.tail = item.prev
	}
}

func (c *lRUCache[K]) Drop() {
	c.Lock()
	defer c.Unlock()

	c.storage = make(map[K]*lruItem[K], c.cap)

	c.head, c.tail = nil, nil
}
