package lru

import "container/list"

// Data Structure It is not safe for concurrent access.
type Cache struct {
	// 键值对 值是指向双向链表的指针
	cache map[string]*list.Element
	maxBytes int64 //最大内存
	nBytes int64 //已用内存
	ll *list.List // 指向双向链表的指针
	// 为什么要有删除后的回调函数？
	OnEvicted func(key string, value Value)
}

type entry struct {
	key string
	// 值可以是任何数据结构 写成一个接口
	value Value
}
type Value interface {
	Len() int
}

// Constructor of Cache
func New(maxBytes int64,onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// implement Len()
func (c *Cache) Len() int{
	return c.ll.Len()
}

// search key's value
func (c *Cache)Get(key string) (value Value,ok bool){
	if node,ok := c.cache[key]; ok{
		c.ll.MoveToFront(node)
		kv := node.Value.(*entry)
		return kv.value,true
	}
	return
}

// remove the oldest
func (c *Cache)RemoveOldest() {
	node := c.ll.Back()
	if node != nil{
		c.ll.Remove(node)
		kv := node.Value.(*entry)
		delete(c.cache,kv.key)
		c.nBytes -= int64(len(kv.key) + kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key,kv.value)
		}
	}
}


// add/modify value
func (c *Cache)Set(key string,value Value) {
	// 如果已经存在key值
	if ele,ok := c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()-kv.value.Len())
		kv.value = value
	} else {
		ele:= c.ll.PushFront(&entry{key,value})
		c.cache[key]=ele
		kv := ele.Value.(*entry)
		// 此处为什么只加一个key？
		c.nBytes += int64(len(key)+kv.value.Len())
	}
	// 将超过最大内存的部分删除 是否需要先删除 把位置腾出来再添加 顺序问题
	for c.maxBytes!=0&&c.nBytes>c.maxBytes{
		c.RemoveOldest()
	}
}
