package caches

import (
	"testing"
	"time"
)

func TestNew_LRU(t *testing.T) {
	cache := New[string](LRU(100))

	switch cache := cache.(type) {
	case *lRUCache[string]:
		if cache.cap != 100 {
			t.Errorf("Expected LRU size 100, got %d", cache.cap)
		}
	default:
		t.Errorf("Expected *lRUCache[int], got %T", cache)
	}
}

func TestNew_ReadOptimized(t *testing.T) {
	cache := New[int](ReadOptimized)

	switch cache := cache.(type) {
	case *readOptimizedCache[int]:
		if cache.timeout != time.Hour {
			t.Errorf("Expected default timeout %v, got %v", time.Hour, cache.timeout)
		}
	default:
		t.Errorf("Expected *readOptimizedCache, got %T", cache)
	}
}

func TestNew_RW(t *testing.T) {
	cache := New[uint8]()

	switch cache := cache.(type) {
	case *Cache[uint8]:
		if cache.timeout != time.Hour {
			t.Errorf("Expected default timeout %v, got %v", time.Hour, cache.timeout)
		}
		if cache.cleanupFreq != time.Hour {
			t.Errorf("Expected default cleanup frequency %v, got %v", time.Hour, cache.cleanupFreq)
		}
	default:
		t.Errorf("Expected *Cache, got %T", cache)
	}
}

func TestNew_LRU_IgnoreTooSmallSize(t *testing.T) {
	cache := New[int](LRU(1))

	switch cache := cache.(type) {
	case *lRUCache[int]:
		if cache.head != nil {
			t.Errorf("Expected default head %v, got %v", nil, cache.head)
		}
		if cache.tail != nil {
			t.Errorf("Expected default  tai %v, got %v", nil, cache.tail)
		}
	default:
		t.Errorf("Expected *lRUCache, got %T", cache)
	}
}

func TestNew_CleanupDisabled(t *testing.T) {
	cache := New[int](CleanupDisabled)

	switch cache := cache.(type) {
	case *Cache[int]:
		if cache.cleanupFreq != time.Hour {
			t.Errorf("Expected default cleanup frequency %v, got %v", time.Hour, cache.cleanupFreq)
		}
	default:
		t.Errorf("Expected *Cache, got %T", cache)
	}
}

func TestWithCleanupFrequency(t *testing.T) {
	cache := New[int](WithCleanupFrequency(time.Minute), WithItemTimeout(time.Hour))

	switch cache := cache.(type) {
	case *Cache[int]:
		if cache.cleanupFreq != 30*time.Second {
			t.Errorf("Expected cleanup frequency 30s, got %v", cache.cleanupFreq)
		}
	default:
		t.Errorf("Expected *Cache, got %T", cache)
	}
}

func TestWithItemTimeout(t *testing.T) {
	cache := New[int](WithItemTimeout(10 * time.Second))

	switch cache := cache.(type) {
	case *Cache[int]:
		if cache.timeout != 10*time.Second {
			t.Errorf("Expected timeout 10s, got %v", cache.timeout)
		}
	default:
		t.Errorf("Expected *Cache, got %T", cache)
	}
}
