package cachehttp

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sancache/ByteCont"
	"sort"
	"testing"
)

var db = map[string]string{
	"JDG": "369",
	"RNG": "888",
	"TES": "000",
}

func Test_ServerHttp(t *testing.T) {
	ByteCont.NewGroup("g1", 1024, ByteCont.GetterFunc(func(key string) ([]byte, error) {
		fmt.Println("group search")
		if con, flag := db[key]; flag {
			return []byte(con), nil
		}
		return nil, errors.New("not found")
	}))

	addr := "127.0.0.1:33366"
	hpool := NewHttpPool(addr)
	//hpool.ServerHTTP()
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, hpool))
}

func Test_SortSearch(t *testing.T) {
	aa := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	gg := sort.Search(len(aa), func(i int) bool {
		if aa[i] > 4 {
			return true
		}
		return false
	})
	fmt.Println(aa, gg, aa[gg])
}
