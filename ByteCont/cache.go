package ByteCont

import (
	"sancache/ByteCont/lru"
	"sync"
)

type cache struct {
	mu           sync.Mutex
	lru          *lru.Cache
	MaxcacheByte int64
}

func (c *cache) add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.NewCache(c.MaxcacheByte, nil)
	}
	c.lru.Add(key, &val)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return ByteView{}, false
	}
	con, flag := c.lru.Get(key)
	if flag == false {
		return ByteView{}, flag
	}
	return *con.(*ByteView), true
}
