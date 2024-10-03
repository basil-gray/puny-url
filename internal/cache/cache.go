package cache

import (
	"fmt"
	"sync"
	"time"

	"puny-url/internal/logger"
)

type Cache struct {
	entries       sync.Map // key: shortID, value: cacheEntry
	evictInterval time.Duration
	ttl           time.Duration
}

type cacheEntry struct {
	LongURL      string
	lastAccessed time.Time
}

func New(evictInterval, ttl time.Duration) *Cache {
	c := &Cache{
		evictInterval: evictInterval,
		ttl:           ttl,
	}
	go c.startEviction()
	return c
}

func (c *Cache) startEviction() {
	for {
		time.Sleep(c.evictInterval)
		now := time.Now()

		c.entries.Range(func(key, value interface{}) bool {
			entry := value.(cacheEntry)
			if now.Sub(entry.lastAccessed) > c.ttl {
				logger.Info(fmt.Sprintf("(Cache expiry) %s: %s has been evicted from the cache", key, entry.LongURL))
				c.entries.Delete(key)
			}
			return true
		})
	}
}

// O(n) in the worst case... can we speed this up??
func (c *Cache) FindByLong(longURL string) (string, bool) {
	var cachedShortID string
	c.entries.Range(func(key, value interface{}) bool {
		entry := value.(cacheEntry)
		if entry.LongURL == longURL {
			cachedShortID = key.(string)
			return false
		}
		return true
	})
	return cachedShortID, cachedShortID != ""
}

func (c *Cache) UpdateCache(shortId string, longURL string) {
	logger.Info(fmt.Sprintf("(Cache update) %s: %s", shortId, longURL))
	c.store(shortId, longURL)
}

func (c *Cache) LoadFromCache(shortId string) (string, bool) {
	longURL, found := c.load(shortId)
	c.store(shortId, longURL) // Every time we touch the cache, refresh
	return longURL, found
}

func (c *Cache) store(shortID string, longURL string) {
	c.entries.Store(shortID, cacheEntry{LongURL: longURL, lastAccessed: time.Now()})
}

func (c *Cache) load(shortID string) (string, bool) {
	entry, found := c.entries.Load(shortID)
	if !found {
		return "", false
	}
	cacheEntry := entry.(cacheEntry)
	return cacheEntry.LongURL, true
}
