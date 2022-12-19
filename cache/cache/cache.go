package cache

import (
	"cache/cache/evict"
	"sync"
)

// 实例化LRU 封装get和add方法 添加互斥锁

type cache struct {
	mu         sync.Mutex // 互斥锁
	evict      *evict.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 延迟初始化
	if c.evict == nil {
		c.evict = evict.New(c.cacheBytes, nil)
	}
	c.evict.Add(key, value)
}

// 从缓存中获取数据
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.evict == nil {
		return
	}
	if v, ok := c.evict.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

// Getter 回调: 从数据源获取数据并添加到缓存
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
