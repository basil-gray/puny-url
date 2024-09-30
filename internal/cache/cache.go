package cache

import (
	"sync"
	"time"

	"puny-url/internal/logger"
)

type Cache struct {
	entries       sync.Map
	evictInterval time.Duration
	ttl           time.Duration
}

type cacheEntry struct {
	LongURL      string
	LastAccessed time.Time
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
			if now.Sub(entry.LastAccessed) > c.ttl {
				logger.Info(key.(string) + ": " + entry.LongURL + " has been evicted from the cache due to inactivity.")
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

func (c *Cache) GetLong(shortID string) (string, bool) {
	entry, found := c.entries.Load(shortID)
	if !found {
		return "", false
	}
	cacheEntry := entry.(cacheEntry)
	return cacheEntry.LongURL, true
}

func (c *Cache) Update(shortID, longURL string) {
	c.entries.Store(shortID, cacheEntry{LongURL: longURL, LastAccessed: time.Now()})
}
