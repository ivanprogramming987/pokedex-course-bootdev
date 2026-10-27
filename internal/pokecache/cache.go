package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
	marked    bool
}

func NewCache(d time.Duration) Cache {
	c := Cache{
		entries: map[string]cacheEntry{},
		mu:      sync.Mutex{},
	}
	go c.reapLoop(d)
	return c
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c Cache) reapLoop(d time.Duration) {
	t := time.NewTicker(d)

	for range t.C {
		c.mu.Lock()
		for entryKey, entry := range c.entries {
			if time.Now().Sub(entry.createdAt) >= d {
				delete(c.entries, entryKey)
			}
		}
		c.mu.Unlock()
	}
}
