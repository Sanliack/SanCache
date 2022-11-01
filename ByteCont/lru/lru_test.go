package lru

import (
	"fmt"
	"testing"
)

type Tdata string

func (t Tdata) Len() int {
	return len(t)
}

func Test_Get(t *testing.T) {
	cache := NewCache(20, func(key string, val Value) {
		fmt.Printf("数据%s:%s已经被删除", key, val)
	})

	cache.Add("key1", Tdata("val1"))
	cache.Add("key2", Tdata("val2"))
	data, flag := cache.Get("key1")
	if flag {
		fmt.Println(data)
	}
}

func Test_lru(t *testing.T) {
	cache := NewCache(20, func(key string, val Value) {
		fmt.Printf("数据%s:%s已经被删除", key, val)
	})

	cache.Add("key1", Tdata("val1"))
	cache.Add("key2", Tdata("val2"))
	cache.Add("key3", Tdata("val3"))
	cache.Add("key4", Tdata("val4"))
	cache.Add("key5", Tdata("val5"))
	cache.Add("key6", Tdata("val6"))
	cache.Add("key7", Tdata("val7"))
	cache.Add("key8", Tdata("val8"))

}
