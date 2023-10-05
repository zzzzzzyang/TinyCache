package TinyCache

import (
	"TinyCache/lru"
	"sync"
)

type cache struct {
	lru        *lru.Cache
	mu         sync.Mutex
	cacheBytes uint32
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if ele, ok := c.lru.Get(key); ok {
		return ele.(ByteView), true
	}
	return
}
