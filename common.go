package caches

import (
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("no value found for privided key")
	ErrKeyExists       = errors.New("key already exists")
	ErrIncompatibleOpt = errors.New("incompatible options passed")
)

type (
	Storage[K comparable] interface {
		Get(key K) (any, error)
		Put(key K, value any)
		Delete(key K)
		Drop()
	}

	Opts func(*cacheOpts)

	cacheOpts struct {
		timeout              time.Duration
		cleanupFreq          time.Duration
		isReadOptimizedCache bool
		cleanupEnabled       bool
		isLRU                bool
		lruSize              int
	}

	item struct {
		val   any
		until time.Time
	}

	lruItem[K comparable] struct {
		key   K
		value any
		prev  *lruItem[K]
		next  *lruItem[K]
	}
)
