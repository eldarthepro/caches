## Package caches
Contains three types of concurrent safe cache for different scenarios. 


## Usage:

```
cache := caches.New[string]()
```
Will instatiate cache with string as key value and default settings: map+mutex based cache with cleanup every 15 minutes with key-value lifetime of 1 hour. Options: WithCleanupFrequency, WithItemTimeout.

```
cache :=  caches.New[int](WithCleanupFrequency(time.Minute), WithItemTimeout(time.Hour))
```
If you need explicitly set desired cleanup and value lifetime.

```
cache := caches.New[int](ReadOptimized)
```
Will instantiate sync.map based cache, good for single write and many reads. Options: WithCleanupFrequency, WithItemTimeout.

```
cache :=  caches.New[string](LRU(100))
```

Will instanitate least recently used cache with given capacity. 

Notice: if both ReadOptimized and LRU() options were passed to constructor it will return lru cache.