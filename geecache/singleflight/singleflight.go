package singleflight

import "sync"

// 请求类型
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// map key和call 映射
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// 并发只执行一次的Do方法
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err

}
