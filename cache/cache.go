package cache

import (
	"log"
	"sync"
	"time"
)

type CacheEntry struct {
	LongURL      string
	Timestamp    time.Time
	LastAccessed time.Time
}

var (
	Map           = sync.Map{}
	EvictInterval time.Duration
	TTL           time.Duration
)

func StartCacheEviction() {
	for {
		time.Sleep(EvictInterval)
		now := time.Now()

		Map.Range(func(key, value interface{}) bool {
			entry := value.(CacheEntry)
			if now.Sub(entry.LastAccessed) > TTL {
				log.Printf("%s: %s has been evicted from the cache due to inactivity.", key, entry.LongURL)
				Map.Delete(key)
			}
			return true
		})
	}
}

func Find(longURL string, now time.Time) (string, bool) {
	var cachedShortID string
	Map.Range(func(key, value interface{}) bool {
		entry := value.(CacheEntry)
		if entry.LongURL == longURL {
			cachedShortID = key.(string)
			Map.Store(cachedShortID, CacheEntry{LongURL: longURL, Timestamp: now, LastAccessed: now})
			return false
		}
		return true
	})
	if cachedShortID != "" {
		return cachedShortID, true
	}
	return "", false
}

func Update(shortID, longURL string, now time.Time) {
	Map.Store(shortID, CacheEntry{LongURL: longURL, Timestamp: now, LastAccessed: now})
}
