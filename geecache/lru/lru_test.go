package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Get(t *testing.T) {
	// maxBytes设为0时代表不对内存设限 不是内存为0
	lru := New(int64(0),nil)
	lru.Set("key1", String("123456"))
	if v,ok:=lru.Get("key1");!ok||string(v.(String))!="123456"{
		t.Fatalf("cache hit key1=1234 failed")
	}
	//if _,ok:=lru.Get("key2");!ok{
	//	t.Fatalf("cache miss key2 failed")
	//}
}
func TestCache_RemoveOldest(t *testing.T) {
	k1,k2,k3 := "Key1","key2","key3"
	v1,v2,v3 := String("value1"), String("value2"), String("value3")
	cap := len(k1+k2)+v1.Len()+v2.Len()
	// 此处keys定义在函数外 传入的是什么？
	keys :=make([]string,0)
	callback := func(key string,value Value) {
		keys=append(keys,key)
	}
	lru:= New(int64(cap),callback)
	lru.Set(k1,v1)
	lru.Set(k2,v2)
	lru.Set(k3,v3)
	if _,ok:=lru.Get("k1");ok {
		t.Fatalf("Removeoldest key1 failed")
	}
	expect:=[]string{"k1"}
	if !reflect.DeepEqual(expect,keys) {
		t.Fatalf("Call failed")
	}
}
