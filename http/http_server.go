package cachehttp

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"sancache/ByteCont"
	"sancache/consistentHash"
	pd "sancache/grpc"
	"strings"
	"sync"
)

var (
	DefaultBasePath = "/_SanCache/"
	defaultReplicas = 50
)

type HttpPool struct {
	self        string
	basepath    string
	mu          sync.Mutex
	peers       *consistentHash.Hashmap
	httpGetters map[string]*httpGetter
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{self: self, basepath: DefaultBasePath}
}

func (h *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[server %s] log: %s", h.self, fmt.Sprintf(format, v...))
}

func (h *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, DefaultBasePath) {
		fmt.Println("http prefix error " + r.URL.Path)
		return
		//panic("http prefix error " + r.URL.Path)
	}
	h.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(h.basepath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := ByteCont.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	body, err := proto.Marshal(&pd.Response{Value: view.GetSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(body)
}

func (h *HttpPool) SetPeers(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.peers = consistentHash.NewHash(defaultReplicas, nil)
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*httpGetter)
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{
			baseURL: peer + h.basepath}
	}
}

func (h *HttpPool) PickPeer(key string) (ByteCont.PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		h.Log("pick peer %s", peer)
		return h.httpGetters[peer], true
	}
	return nil, false
}
