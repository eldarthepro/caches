package caches

import (
	"testing"
)

func TestLRUCache_PutAndGet(t *testing.T) {
	cache := newLRU[int](cacheOpts{lruSize: 3})

	cache.Put(1, "dog")
	cache.Put(2, "cat")
	cache.Put(3, "mouse")

	if val, err := cache.Get(1); err != nil || val != "dog" {
		t.Errorf("Expected 'one', got %v, error: %v", val, err)
	}
	if val, err := cache.Get(2); err != nil || val != "cat" {
		t.Errorf("Expected 'two', got %v, error: %v", val, err)
	}

	if _, err := cache.Get(4); err == nil {
		t.Error("Expected error for missing key 4, got nil")
	}
}

func TestLRUCache_PutEvictsTailWhenFull(t *testing.T) {
	cache := newLRU[int](cacheOpts{lruSize: 3})

	cache.Put(1, "dog")
	cache.Put(2, "cat")
	cache.Put(3, "mouse")

	cache.Put(4, "fish")

	if _, err := cache.Get(1); err == nil {
		t.Error("Expected error for evicted key 1, got nil")
	}

	if val, err := cache.Get(2); err != nil || val != "cat" {
		t.Errorf("Expected 'two', got %v, error: %v", val, err)
	}
	if val, err := cache.Get(3); err != nil || val != "mouse" {
		t.Errorf("Expected 'three', got %v, error: %v", val, err)
	}
	if val, err := cache.Get(4); err != nil || val != "fish" {
		t.Errorf("Expected 'four', got %v, error: %v", val, err)
	}
}

func TestLRUCache_Delete(t *testing.T) {
	cache := newLRU[int](cacheOpts{lruSize: 3})

	cache.Put(1, "dog")
	cache.Put(2, "cat")

	cache.Delete(1)

	if _, err := cache.Get(1); err == nil {
		t.Error("Expected error for missing key 1, got nil")
	}

	if val, err := cache.Get(2); err != nil || val != "cat" {
		t.Errorf("Expected 'cat', got %v, error: %v", val, err)
	}
}

func TestLRUCache_Drop(t *testing.T) {
	cache := newLRU[int](cacheOpts{lruSize: 3})

	cache.Put(1, "dog")
	cache.Put(2, "cat")
	cache.Put(3, "mouse")

	cache.Drop()

	if _, err := cache.Get(1); err == nil {
		t.Error("Expected error for missing key 1, got nil")
	}
	if _, err := cache.Get(2); err == nil {
		t.Error("Expected error for missing key 2, got nil")
	}
	if _, err := cache.Get(3); err == nil {
		t.Error("Expected error for missing key 3, got nil")
	}
}

func TestLRUCache_MoveToFront(t *testing.T) {
	cache := newLRU[int](cacheOpts{lruSize: 3})

	cache.Put(1, "dog")
	cache.Put(2, "cat")
	cache.Put(3, "mouse")

	entry := cache.storage[1]
	cache.moveToFront(entry)

	if cache.head.key != 1 {
		t.Errorf("Expected head to be key 1, but got %v", cache.head.key)
	}

	if cache.tail.key != 2 {
		t.Errorf("Expected tail to be key 2, but got %v", cache.tail.key)
	}
}
