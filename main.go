package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sancache/ByteCont"
	cachehttp "sancache/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *ByteCont.Group {
	return ByteCont.NewGroup("scores", 2<<10, ByteCont.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *ByteCont.Group) {
	peers := cachehttp.NewHttpPool(addr)
	peers.SetPeers(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *ByteCont.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(view.GetSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

//func main() {
//	aa := []int{4: 2, 3: 4}
//	fmt.Println(aa)
//}
func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}

//
//var ansdb = map[string]string{
//	"key1": "val1",
//	"key2": "val2",
//	"key3": "val3",
//}

//func ls(addr string) {
//	http.ListenAndServe(,)
//}

//func main() {
//	g1 := ByteCont.NewGroup("g1", 1024, ByteCont.GetterFunc(func(key string) ([]byte, error) {
//		fmt.Println("===========")
//		if con, flag := ansdb[key]; flag {
//			return []byte(con), nil
//		}
//		return nil, errors.New("not find")
//	}))
//	peers := cachehttp.NewHttpPool("pool1")
//	addr := []string{"http://localhost:8001", "http://localhost:8002", "http://localhost:8003"}
//	peers.SetPeers(addr...)
//	for _, v := range addr {
//		go ls(v)
//	}
//
//	g1.RegisterPeers(peers)
//	//http.Handle("/sanli", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
//	//	writer.WriteHeader(http.StatusOK)
//	//	writer.Write([]byte("sanliack"))
//	//}))
//	log.Fatal(http.ListenAndServe("localhost:33344", peers))
//
//}
