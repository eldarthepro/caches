package caches

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := newRW[string](cacheOpts{timeout: time.Minute, cleanupFreq: time.Minute})

	cache.Put("key_one", "supman")
	cache.Put("key_two", "badman")

	// Get values from the cache
	if val, err := cache.Get("key_one"); err != nil || val != "supman" {
		t.Errorf("Expected 'supman', got %v, error: %v", val, err)
	}
	if val, err := cache.Get("key_two"); err != nil || val != "badman" {
		t.Errorf("Expected 'badman', got %v, error: %v", val, err)
	}

	// Test Get for a non-existent key
	if _, err := cache.Get("09090"); err == nil {
		t.Error("Expected error for missing key '09090', got nil")
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := newRW[string](cacheOpts{timeout: 1 * time.Millisecond, cleanupFreq: 1 * time.Second})

	cache.Put("key", "value")

	time.Sleep(2 * time.Millisecond)

	if _, err := cache.Get("key"); err == nil {
		t.Error("Expected error for expired key 'key', got nil")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := newRW[int](cacheOpts{timeout: time.Minute, cleanupFreq: time.Minute})

	cache.Put(1, "val1")
	cache.Put(2, "val2")

	cache.Delete(1)

	if _, err := cache.Get(1); err == nil {
		t.Error("Expected error for deleted key '1', got nil")
	}

	if val, err := cache.Get(2); err != nil || val != "val2" {
		t.Errorf("Expected '2', got %v, error: %v", val, err)
	}
}

func TestCache_Drop(t *testing.T) {
	cache := newRW[string](cacheOpts{timeout: 1 * time.Second, cleanupFreq: 1 * time.Second})

	// Set items in the cache
	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	// Drop all items from the cache
	cache.Drop()

	// Test that all items were dropped
	if _, err := cache.Get("key1"); err == nil {
		t.Error("Expected error for missing key 'key1', got nil")
	}
	if _, err := cache.Get("key2"); err == nil {
		t.Error("Expected error for missing key 'key2', got nil")
	}
}


func TestCache_SchedCleanup(t *testing.T) {
	cache := newRW[uint8](cacheOpts{timeout: 1 * time.Millisecond, cleanupFreq: 10 * time.Millisecond})

	cache.Put(1, "value1")
	cache.Put(2, "value2")

	time.Sleep(20 * time.Millisecond)

	if _, err := cache.Get(1); err == nil {
		t.Error("Expected error for expired key '1', got nil")
	}
	if _, err := cache.Get(2); err == nil {
		t.Error("Expected error for expired key '2', got nil")
	}
}
