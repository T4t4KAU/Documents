package flight

import "sync"

// 正在进行中或已经结束的请求
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mutex sync.Mutex
	calls map[string]*call
}

func (group *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	group.mutex.Lock()
	if group.calls == nil {
		group.calls = make(map[string]*call)
	}

	// 检查是否有key的请求 如果有请求则等待并返回
	if c, ok := group.calls[key]; ok {
		group.mutex.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	// 第一次key的请求 记录到表
	c := new(call)
	c.wg.Add(1)
	group.calls[key] = c
	group.mutex.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	group.mutex.Lock()
	delete(group.calls, key)
	group.mutex.Unlock()

	return c.val, c.err
}
