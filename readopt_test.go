package caches

import (
	"testing"
	"time"
)

func TestReadOptimizedCache_SetAndGet(t *testing.T) {
	cache := newReadOptimized[string](cacheOpts{timeout: time.Minute, cleanupFreq: time.Minute})

	cache.Put("key1", "dog")
	cache.Put("key2", "cat")

	if val, err := cache.Get("key1"); err != nil || val != "dog" {
		t.Errorf("Expected 'value1', got %v, error: %v", val, err)
	}
	if val, err := cache.Get("key2"); err != nil || val != "cat" {
		t.Errorf("Expected 'value2', got %v, error: %v", val, err)
	}

	if _, err := cache.Get("key3"); err == nil {
		t.Error("Expected error for missing key 'key3', got nil")
	}
}

func TestReadOptimizedCache_Expiration(t *testing.T) {
	cache := newReadOptimized[int](cacheOpts{timeout: 1 * time.Millisecond, cleanupFreq: 1 * time.Second})

	cache.Put(1, "dog")

	time.Sleep(2 * time.Millisecond)

	if _, err := cache.Get(1); err == nil {
		t.Error("Expected error for expired key '1', got nil")
	}
}

func TestReadOptimizedCache_Delete(t *testing.T) {
	cache := newReadOptimized[string](cacheOpts{timeout: 1 * time.Second, cleanupFreq: 1 * time.Second})

	cache.Put("key1", "spider")
	cache.Put("key2", "wolf")

	cache.Delete("key1")

	if _, err := cache.Get("key1"); err == nil {
		t.Error("Expected error for deleted key 'spider', got nil")
	}

	if val, err := cache.Get("key2"); err != nil || val != "wolf" {
		t.Errorf("Expected 'wolf', got %v, error: %v", val, err)
	}
}

func TestReadOptimizedCache_Drop(t *testing.T) {
	cache := newReadOptimized[string](cacheOpts{timeout: 1 * time.Second, cleanupFreq: 1 * time.Second})

	cache.Put("key1", "spider")
	cache.Put("key2", "wolf")

	cache.Drop()

	if _, err := cache.Get("key1"); err == nil {
		t.Error("Expected error for missing key 'spider', got nil")
	}
	if _, err := cache.Get("key2"); err == nil {
		t.Error("Expected error for missing key 'wolf', got nil")
	}
}

func TestReadOptimizedCache_SchedCleanup(t *testing.T) {
	cache := newReadOptimized[int16](cacheOpts{timeout: 1 * time.Millisecond, cleanupFreq: 10 * time.Millisecond})

	cache.Put(2, "spider")
	cache.Put(4, "wolf")

	time.Sleep(15 * time.Millisecond)

	if _, err := cache.Get(2); err == nil {
		t.Error("Expected error for expired key '2', got nil")
	}
	if _, err := cache.Get(4); err == nil {
		t.Error("Expected error for expired key '4', got nil")
	}
}
