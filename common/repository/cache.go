package repository

import (
	"os"
	"strconv"
	"time"

	"gopkg.in/redis.v5"

	lru "go-web-app/thirdparty/golang-lru"
)

const (
	defaultCacheStoreTTLHours = 24
)

var (
	defaultCacheStoreTTL = defaultCacheStoreTTLHours * time.Hour
)

type CacheStore struct {
	memCache    *lru.Cache
	redisClient *redis.Cmdable
}

func init() {
	cacheStoreTTLStr := os.Getenv("CACHE_STORE_TTL")
	if cacheStoreTTLStr != "" {
		val, err := strconv.ParseFloat(cacheStoreTTLStr, 64)
		if err != nil && val > 0 {
			defaultCacheStoreTTL = time.Duration(val * float64(time.Hour))
		}
	}
}

func NewCacheStoreWithRedisClient(redisClient *redis.Cmdable) *CacheStore {
	memCache, _ := lru.New(50000)
	return &CacheStore{
		memCache:    memCache,
		redisClient: redisClient,
	}
}

// Get looks up a key's value from the cache.
func (c *CacheStore) Get(key string) (interface{}, bool) {
	if c.redisClient != nil {
		exists := (*c.redisClient).Exists(key).Val()
		if !exists {
			return nil, false
		}

		result, err := (*c.redisClient).Get(key).Result()
		if err != nil {
			return nil, false
		}

		return result, true
	}

	return c.memCache.Get(key)
}

/** Set sets a value to the cache.
 */
func (c *CacheStore) Set(key string, value interface{}) interface{} {
	duration := 24 * time.Hour
	return c.SetWithTTL(key, value, &duration)
}

/** SetWithTTL sets a value to the cache. param:ttl could be nil(default ttl)
 */
func (c *CacheStore) SetWithTTL(key string, value interface{}, ttl *time.Duration) interface{} {
	if c.redisClient != nil {

		var ttlVal time.Duration
		if ttl != nil {
			ttlVal = *ttl
		} else {
			ttlVal = 24 * time.Hour
		}

		result, err := (*c.redisClient).Set(key, value, ttlVal).Result()
		if err != nil {
			return nil
		}

		return result
	}

	return c.memCache.Add(key, value)
}

// Del removes the provided key from the cache.
func (c *CacheStore) Del(key string) bool {
	if c.redisClient != nil {
		count, err := (*c.redisClient).Del(key).Result()
		if err != nil {
			return false
		}
		return count > 0
	}

	return c.memCache.Remove(key)
}
