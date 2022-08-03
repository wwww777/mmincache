package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 获取hash值
type Hash func(data []byte) uint32

// 一致性算法的数据结构 hash函数 虚拟节点倍数replicas 哈希环keys 真实节点名与虚拟节点hash值的映射表hashMap
type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

// Map的构造函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 添加真实节点的add方法 入参keys为多个节点的名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 计算新加节点的虚拟节点hash值
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// key选择节点的get方法
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// 找到第一个大于key的哈希值的所在的索引
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 返回真实节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
