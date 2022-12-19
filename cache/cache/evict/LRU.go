package evict

import (
	"container/list"
)

// LRU算法: 最近最少使用，如果数据最近被访问过，那么将来被访问的概率也会更高
// 维护一个队列，则移动到队尾，那么队首则是最近最少访问的数据，淘汰该记录

type Cache struct {
	maxBytes  int64                         // 最大内存
	nBytes    int64                         // 当前已使用内存
	List      *list.List                    // 双向链表
	cache     map[string]*list.Element      // 字典
	OnEvicted func(key string, value Value) // 回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New 实例化
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		List:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找: 从字典中找到对应的双向链表的节点 将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.List.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 删除: 淘汰缓存 移除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	element := c.List.Back() // 队首元素
	if element != nil {
		c.List.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key) // 在字典中删除
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 增加/修改: 如果键存在 则更新对应的节点 将该节点移动到队尾
func (c *Cache) Add(key string, value Value) {
	// 如果键存在 则更新对应节点的值 将该节点移动到队尾
	// 不存在则新增 在队尾添加新节点 并在字典中添加KV
	if element, ok := c.cache[key]; ok {
		c.List.MoveToFront(element)
		kv := element.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
	} else {
		element := c.List.PushFront(&entry{key, value})
		c.cache[key] = element
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	// 如果超过了设定的最大值 则移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.List.Len()
}
