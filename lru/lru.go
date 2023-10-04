package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nBytes    int64
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

func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
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
		c.nBytes -= int64(len(kv.key)) + int64(kv.val.Len())
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
		c.nBytes += int64(kv.val.Len()) - int64(value.Len())
	} else {
		ele = c.ll.PushFront(&entry{key: key, val: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
