package cache

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	Entries       sync.Map
	evictInterval time.Duration
	ttl           time.Duration
}

type cacheEntry struct {
	LongURL      string
	Timestamp    time.Time
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

		c.Entries.Range(func(key, value interface{}) bool {
			entry := value.(cacheEntry)
			if now.Sub(entry.LastAccessed) > c.ttl {
				log.Printf("%s: %s has been evicted from the cache due to inactivity.", key, entry.LongURL)
				c.Entries.Delete(key)
			}
			return true
		})
	}
}

func (c *Cache) FindByLong(longURL string) (string, bool) {
	var cachedShortID string
	c.Entries.Range(func(key, value interface{}) bool {
		entry := value.(cacheEntry)
		if entry.LongURL == longURL {
			cachedShortID = key.(string)
			c.Update(cachedShortID, longURL)
			return false
		}
		return true
	})
	return cachedShortID, cachedShortID != ""
}

func (c *Cache) GetLong(shortURL string) (string, bool) {
	entry, found := c.Entries.Load(shortURL)
	if !found {
		return "", false
	}
	cacheEntry := entry.(cacheEntry)
	return cacheEntry.LongURL, true
}

func (c *Cache) Update(shortID, longURL string) {
	now := time.Now()
	c.Entries.Store(shortID, cacheEntry{LongURL: longURL, Timestamp: now, LastAccessed: now})
}
