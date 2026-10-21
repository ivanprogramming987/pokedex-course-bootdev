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

func NewCache(time time.Duration) Cache {
	c := Cache{
		entries: map[string]cacheEntry{},
		mu:      sync.Mutex{},
	}
	go c.reapLoop(time)
	return c
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] := cacheEntry{
		createdAt: time.Now(),
		val:       val,
		marked:    false,
	}
}

func (c Cache) Get(key string) []byte, bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.entries[key]
	return v, ok
}

func (c Cache) reapLoop(time time.Duration) {
	t := time.NewTicker(time)

	for range t.C {
		c.mu.Lock()
		for entryKey, entry := range c.entries {
			if entry.marked == true {
				delete(c.entries, entryKey)
			} else {
				entry.marked = true
			}
		}
		c.mu.Unlock()
	}
}
