package cache

import (
	"log"
	"sync"
	"time"
)

type CacheEntry struct {
	Value        string
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
				log.Printf("%s: %s has been evicted from the cache due to inactivity.", key, entry.Value)
				Map.Delete(key)
			}
			return true
		})
	}
}
