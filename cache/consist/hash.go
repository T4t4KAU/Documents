package consist

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 一致性哈希: 将key映射到2^32的空间中 将数字首位相连 形成一个环
// 计算节点/机器(通常使用节点的名称、编号和IP地址)的哈希值 放置在环上
// 计算key的哈希值 放置在环上 顺时针寻找到的一个节点 就是应选取的节点/机器
// 一致性哈希算法在新增/删除节点时 只要重新定位该节点附近的一小部分数据 而无需重新定位所有的节点

// 数据倾斜: 如果服务器节点过少 容易引起key的倾斜 最终使得缓存节点负载不均
// 引入虚拟节点 一个真实节点对应多个虚拟节点

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

// New 创建一个Map实例
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	// 哈希算法默认采用checksumIEEE
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 增加新节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 创建replicas个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 通过添加编号的方式区分不同虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash) // 添加到环上
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 选择节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key))) // 计算哈希值

	// 获取第一个匹配的虚拟节点下标
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[index%len(m.keys)]]
}

// Sub 删除节点
func (m *Map) Sub(keys ...string) {

}
