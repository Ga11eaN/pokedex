package pokecache


import (
    "time"
    "sync"
)


type Cache struct {
    cache map[string]cacheEntry
    mu sync.Mutex
}

type cacheEntry struct {
    createdAt time.Time
    val []byte
}

func NewCache(interval time.Duration) Cache {
    c := Cache{
        cache: make(map[string]cacheEntry),
        mu: sync.Mutex{},
    }

    go c.reapLoop(interval)

    return c
}

func (c *Cache) reapLoop(interval time.Duration) {
    ticker := time.NewTicker(interval / 3)
    for range ticker.C {
        expiredTime := time.Now().Add(-interval)
        c.mu.Lock()
        for k, v := range c.cache {
            if v.createdAt.Before(expiredTime) {
                delete(c.cache, k)
            }
        }
        c.mu.Unlock()
    }
}

func (c *Cache) Add(key string, val []byte) {
    cNew := cacheEntry {
        createdAt: time.Now(),
        val: val,
    }
    c.mu.Lock()
    c.cache[key] = cNew
    c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    value, ok := c.cache[key]
    c.mu.Unlock()
    if !ok {
        return nil, false
    }
    return value.val, true
}