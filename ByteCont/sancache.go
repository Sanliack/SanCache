package ByteCont

import (
	"errors"
	pd "sancache/grpc"
	"sancache/singleflight"
	"sync"
)

var (
	mu       sync.RWMutex
	groupmap = make(map[string]*Group)
)

// 接口式函数
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

type Group struct {
	name   string
	getter Getter
	cache  cache
	peers  PeerPicker
	loader *singleflight.SFGroup
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key can not null")
	}
	bview, flag := g.cache.get(key)
	if flag {
		return bview, nil
	}
	return g.load(key)
}

//func (g *Group) load(key string) (ByteView, error) {
//	return g.getLocally(key)
//}

func (g *Group) getLocally(key string) (ByteView, error) {
	bview, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	val := ByteView{b: bview}
	g.populateCache(key, val)
	return val, nil
}

func (g *Group) populateCache(key string, val ByteView) {
	g.cache.add(key, val)
}

func (g *Group) RegisterPeers(peer PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peer

}

func (g *Group) load(key string) (ByteView, error) {
	bview, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, flag := g.peers.PickPeer(key); flag {
				bv, err := g.GetFromPeer(peer, key)
				if err != nil {
					return bv, err
				}
			}
		}
		return g.getLocally(key)
	})
	if err != nil {
		return ByteView{}, err
	}
	return bview.(ByteView), nil

}

func (g *Group) GetFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pd.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pd.Response{}
	bytes, err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func NewGroup(name string, MaxcacheByte int64, getter Getter) *Group {
	if getter == nil {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	newg := &Group{
		name:   name,
		getter: getter,
		cache:  cache{MaxcacheByte: MaxcacheByte},
		loader: &singleflight.SFGroup{},
	}
	groupmap[name] = newg
	return newg
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groupmap[name]
}
