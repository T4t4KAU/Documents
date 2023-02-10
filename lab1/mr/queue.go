package mr

import "sync"

type queue struct {
	data []string
	mux  sync.RWMutex
}

func Queue() *queue {
	var q queue
	q.data = make([]string, 0)
	q.mux = sync.RWMutex{}
	return &q
}

func (q *queue) Enqueue(val string) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.data = append(q.data, val)
}

func (q *queue) Dequeue() string {
	q.mux.Lock()
	defer q.mux.Unlock()
	item := q.data[0]
	q.data = q.data[1:]
	return item
}

func (q *queue) isEmpty() bool {
	return len(q.data) == 0
}
