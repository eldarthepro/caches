package caches

import (
	"time"
)

// New returns concurrent-safe cache implementation depending on options provided. By default returns map based cache that performs well im most scenarios.
func New[K comparable](opts ...Opts) Storage[K] {
	o := defaultOpts()
	for _, opt := range opts {
		opt(&o)
	}

	if o.isLRU {
		return newLRU[K](o)
	}

	if o.isReadOptimizedCache {
		return newReadOptimized[K](o)
	}

	return newRW[K](o)
}

func defaultOpts() cacheOpts {
	return cacheOpts{
		timeout:        time.Hour,
		cleanupFreq:    time.Minute * 15,
		cleanupEnabled: true,
		lruSize:        500,
	}
}

// LRU use if least recently used cache is needed, if size is smaller than 1 will set size to 500, not will ignore option ReadOptimized if passed
func LRU(size int) Opts {
	return func(o *cacheOpts) {
		o.isLRU = true
		o.lruSize = size

		if o.lruSize < 2 {
			o.lruSize = 500
		}
	}
}

// ReadOptimized use if value for given key is not expected to change over time, will be ignored if LRU option is passed
func ReadOptimized(o *cacheOpts) {
	o.isReadOptimizedCache = true
}

// CleanupDisabled allows to disable cache cleanup
func CleanupDisabled(o *cacheOpts) {
	o.cleanupEnabled = false
}

// WithCleanupFrequency allows to set custom cleanup frequency
func WithCleanupFrequency(freq time.Duration) Opts {
	return func(o *cacheOpts) {
		o.cleanupFreq = freq
	}
}

// WithItemTimeout allows to set custom item expiry
func WithItemTimeout(timeout time.Duration) Opts {
	return func(o *cacheOpts) {
		o.timeout = timeout
	}
}

// NeverExpire returns maximum possible duration
func NeverExpire() time.Duration {
	return time.Duration(1<<63 - 1)
}

// DefaultExpireDuration returns duration of one hour
func DefaultExpireDuration() time.Duration {
	return time.Hour
}
