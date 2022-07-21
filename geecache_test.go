package geecache

import (
	"fmt"
	"log"
	"testing"
)

// 模拟数据库
var db = map[string]string{
	"Tom":"630",
	"Jack": "589",
	"Sam":  "567",
}

// 测试Get方法
func TestGet(t *testing.T)  {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key",key)
		if v,ok := db[key];ok{
			if _,ok:=loadCounts[key];!ok {
				loadCounts[key]=0
			}
			loadCounts[key]+=1;
			return []byte(v),nil
		}
		return nil,fmt.Errorf("%s not exist",key)
	}),2<<10)
	for k,v :=range db {
		if view,err:=gee.Get(k);err!=nil||view.String()!=v{
			t.Fatal("failed to get value of Tom")
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}