package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu sync.Mutex
	lru *lru.Cache
	//
	cacheBytes int64
}

func (c *cache)add(key string,value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 延迟初始化
	if c.lru == nil{
		c.lru = lru.New(c.cacheBytes,nil)
	}
	c.lru.Set(key,value)
}

func (c *cache)get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil{
		// 不用报错吗？
		return
	}

	if v,ok :=c.lru.Get(key);ok{
		// 类型强制转换？
		return v.(ByteView),true
	}
	return
}

