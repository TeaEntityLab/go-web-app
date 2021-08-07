package lru

import (
	"time"
)

// CacheWithExpiration is a thread-safe fixed size LRU cache.
type CacheWithExpiration struct {
	lru            *Cache
	expirationTime time.Duration
}

type CacheWithExpirationEntity struct {
	timestamp    time.Time
	value        interface{}
	neverExpired bool
}

// New creates an LRU of the given size.
func NewCacheWithExpiration(size int, expirationTime time.Duration) (*CacheWithExpiration, error) {
	return NewCacheWithExpirationAndEvict(size, expirationTime, nil)
}

// NewCacheWithExpirationAndEvict constructs a fixed size cache with the given eviction
// callback.
func NewCacheWithExpirationAndEvict(size int, expirationTime time.Duration, onEvicted func(key interface{}, value interface{})) (*CacheWithExpiration, error) {
	lru, err := NewWithEvict(size, onEvicted)
	if err != nil {
		return nil, err
	}
	c := &CacheWithExpiration{
		lru:            lru,
		expirationTime: expirationTime,
	}
	return c, nil
}

// Purge is used to completely clear the cache.
func (c *CacheWithExpiration) Purge() {
	c.lru.Purge()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *CacheWithExpiration) Add(key, value interface{}) (evicted bool) {
	evicted = c.lru.Add(key, CacheWithExpirationEntity{
		timestamp: time.Now(),
		value:     value,
	})
	return evicted
}

// Get looks up a key's value from the cache.
func (c *CacheWithExpiration) Get(key interface{}) (value interface{}, ok bool) {
	value, ok = c.lru.Get(key)
	if ok && value != nil {
		value, ok = c.checkNotTimeoutByRawKeyValue(key, value)
		if !ok {
			value = nil
		}
	}
	return value, ok
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
func (c *CacheWithExpiration) Contains(key interface{}) bool {
	value, ok := c.Peek(key)
	if ok && value != nil {
		value, ok = c.checkNotTimeoutByRawKeyValue(key, value)
		if !ok {
			return false
		}
	}
	containKey := c.lru.Contains(key)
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *CacheWithExpiration) Peek(key interface{}) (value interface{}, ok bool) {
	value, ok = c.lru.Peek(key)
	if ok && value != nil {
		value, ok = c.checkNotTimeoutByRawKeyValue(key, value)
		if !ok {
			value = nil
		}
	}
	return value, ok
}

// ContainsOrAdd checks if a key is in the cache  without updating the
// recent-ness or deleting it for being stale,  and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *CacheWithExpiration) ContainsOrAdd(key, value interface{}) (ok, evicted bool) {
	if c.Contains(key) {
		return true, false
	}
	evicted = c.Add(key, value)
	return false, evicted
}

// Remove removes the provided key from the cache.
func (c *CacheWithExpiration) Remove(key interface{}) {
	c.lru.Remove(key)
}

// RemoveOldest removes the oldest item from the cache.
func (c *CacheWithExpiration) RemoveOldest() {
	c.lru.RemoveOldest()
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *CacheWithExpiration) Keys() []interface{} {
	keys := c.lru.Keys()
	return keys
}

// Len returns the number of items in the cache.
func (c *CacheWithExpiration) Len() int {
	length := c.lru.Len()
	return length
}

func (c *CacheWithExpiration) checkNotTimeoutByRawKeyValue(key interface{}, value interface{}) (interface{}, bool) {
	entry := value.(CacheWithExpirationEntity)
	if !entry.neverExpired && time.Now().Sub(entry.timestamp) > c.expirationTime {
		c.Remove(key)
		return nil, false
	}

	return entry.value, true
}
