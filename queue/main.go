package queue

import (
	"sync"
)

type Qu struct {
	arr []interface{}
	mu  *sync.RWMutex
}

func (q *Qu) Unshift(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	newArr := make([]interface{}, len(q.arr)+1)
	newArr[0] = item
	for i := 0; i < len(q.arr); i++ {
		newArr[i+1] = q.arr[i]
	}
	q.arr = newArr
}

func (q *Qu) Push(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.arr = append(q.arr, item)
}

func (q *Qu) Pop() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.arr) == 0 {
		return nil
	}
	item := q.arr[0]
	q.arr = q.arr[1:len(q.arr)]
	return item
}

func (q Qu) Len() int {
	return len(q.arr)
}

func (q *Qu) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.arr = []interface{}{}
}

func New() *Qu {
	q := new(Qu)
	q.mu = &sync.RWMutex{}
	q.arr = make([]interface{}, 0)
	return q
}
