package mr

import "sync"

type item struct {
	task int32
	file string
}

type queue struct {
	data []item
	mux  sync.RWMutex
}

func Queue() *queue {
	var q queue
	q.data = make([]item, 0)
	q.mux = sync.RWMutex{}
	return &q
}

func (q *queue) Enqueue(i item) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.data = append(q.data, i)
}

func (q *queue) Dequeue() item {
	q.mux.Lock()
	defer q.mux.Unlock()
	item := q.data[0]
	q.data = q.data[1:]
	return item
}

func (q *queue) empty() bool {
	return len(q.data) == 0
}
