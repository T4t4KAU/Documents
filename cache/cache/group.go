package cache

import (
	"cache/clog"
	"cache/flight"
	"cache/service/peers"
	"fmt"
	"log"
	"sync"
)

// Group 负责与用户交互 控制缓存值存储和获取流程
type Group struct {
	name      string // 唯一名称
	getter    Getter // 回调获取源数据
	mainCache cache  // 并发缓存
	peers     peers.PeerPicker
	loader    *flight.Group
}

var (
	mutex  sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 实例化Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mutex.Lock()
	defer mutex.Unlock()
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &flight.Group{},
	}
	groups[name] = group
	return group
}

// GetGroup 返回一个命名的group
func GetGroup(name string) *Group {
	mutex.RLock()
	group := groups[name]
	mutex.RUnlock()
	return group
}

// Get 从mainCache中查找缓存
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 如果存在则返回缓存值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Cache] hit")
		return v, nil
	}
	return g.load(key)
}

// RegisterPeers 注册节点
func (g *Group) RegisterPeers(peers peers.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// 缓存缺失时加载数据
func (g *Group) load(key string) (value ByteView, err error) {
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				clog.Info("Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return view.(ByteView), nil
	}
	return
}

// 本地获取数据
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	return value, nil
}

// 从其他节点获取数据
func (g *Group) getFromPeer(peer peers.PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
