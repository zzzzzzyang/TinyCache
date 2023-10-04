package lru

import "container/list"

// todo 添加LFU和LRU-K

type Cache struct {
	maxBytes  uint32
	nBytes    uint32
	cache     map[string]*list.Element
	ll        *list.List
	OnEvicted func(key string, value Value)
}

type entry struct {
	key string
	val Value
}

type Value interface {
	Len() int
}

func NewCache(maxBytes uint32, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.val, true
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= uint32(len(kv.key)) + uint32(kv.val.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.val)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		kv.val = value
		c.nBytes += uint32(kv.val.Len()) - uint32(value.Len())
	} else {
		ele = c.ll.PushFront(&entry{key: key, val: value})
		c.cache[key] = ele
		c.nBytes += uint32(len(key)) + uint32(value.Len())
	}
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
