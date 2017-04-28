package queue

import (
	"container/heap"
	"sync"
)

type Interface interface {
	Less(other interface{}) bool
}

type sorter []Interface

// Implement heap.Interface: Push, Pop, Len, Less, Swap
func (s *sorter) Push(x interface{}) {
	*s = append(*s, x.(Interface))
}

func (s *sorter) Pop() interface{} {
	n := len(*s)
	if n > 0 {
		x := (*s)[n-1]
		*s = (*s)[0 : n-1]
		return x
	}
	return nil
}

func (s *sorter) Len() int {
	return len(*s)
}

func (s *sorter) Less(i, j int) bool {
	return (*s)[i].Less((*s)[j])
}

func (s *sorter) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Define priority queue struct
type PriorityQueue struct {
	s       *sorter
	rwMutex *sync.RWMutex
}

func NewPriorityQueue() *PriorityQueue {
	q := &PriorityQueue{
		s:       new(sorter),
		rwMutex: &sync.RWMutex{},
	}

	heap.Init(q.s)
	return q
}

func (q *PriorityQueue) Push(x Interface) {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()
	heap.Push(q.s, x)
}

func (q *PriorityQueue) Pop() Interface {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()
	if len(*q.s) > 0 {
		return heap.Pop(q.s).(Interface)
	}
	return nil
}

func (q *PriorityQueue) Top() Interface {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()
	if len(*q.s) > 0 {
		return (*q.s)[0].(Interface)
	}
	return nil
}

func (q *PriorityQueue) Fix(x Interface, i int) {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()
	(*q.s)[i] = x
	heap.Fix(q.s, i)
}

func (q *PriorityQueue) Remove(i int) Interface {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()
	return heap.Remove(q.s, i).(Interface)
}

func (q *PriorityQueue) Len() int {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()
	return q.s.Len()
}
