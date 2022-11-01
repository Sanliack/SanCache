package lru

import "container/list"

type Cache struct {
	//最大使用空间
	maxbytes int64
	//已使用空间
	nbytes int64
	ll     *list.List
	cache  map[string]*list.Element
	//被移除时的回调函数
	OnEvicted func(key string, val Value)
}

type entry struct {
	key string
	val Value
}

type Value interface {
	Len() int
}

func NewCache(maxbytes int64, onevicted func(key string, val Value)) *Cache {
	return &Cache{
		maxbytes:  maxbytes,
		nbytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onevicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	if tar, flag := c.cache[key]; flag {
		c.ll.MoveToFront(tar)
		kv := tar.Value.(*entry)
		return kv.val, true
	}
	return nil, false
}

func (c *Cache) RemoveOldestCache() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.val.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.val)
		}
	}
}

func (c *Cache) Add(key string, val Value) {
	if ele, flag := c.cache[key]; flag {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(val.Len()) - int64(kv.val.Len())
		kv.val = val
	} else {
		ent := &entry{
			val: val,
			key: key,
		}
		ele := c.ll.PushFront(ent)
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(val.Len())
	}
	for c.maxbytes != 0 && c.maxbytes < c.nbytes {
		c.RemoveOldestCache()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
