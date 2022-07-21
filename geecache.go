package geecache

import (
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte,error)
}

type GetterFunc func(key string)([]byte,error)

func (f GetterFunc)Get(key string) ([]byte,error){
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

// 构造函数
func NewGroup(name string, getter Getter,cacheBytes int64) *Group {
	if getter == nil{
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g :=&Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// 从groups中获取group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get方法
func (g *Group)Get(key string) (ByteView,error) {
	if v,ok:=g.mainCache.get(key);ok{
		log.Println("[Cache] hit")
		return v,nil
	}
	return g.load(key)
}

func (g *Group)load(key string) (ByteView,error){
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView,error){
	byte,err := g.getter.Get(key)
	if err != nil{
		return ByteView{},err
	}
	// 获取拷贝值返回
	value := ByteView{b: cloneBytes(byte)}
	// 放入缓存中
	g.populateCache(key, value)
	return value,nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key,value)
}
