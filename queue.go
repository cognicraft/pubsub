package pubsub

import (
	"fmt"
	"sort"
	"sync"
)

type Queue interface {
	Enqueue(m Message) error
	Dequeue() (Message, error)
	Peek() (Message, error)
	IsEmpty() bool
}

func NewSimpleQueue() *SimpleQueue {
	return &SimpleQueue{}
}

type SimpleQueue struct {
	mu       sync.RWMutex
	messages []Message
}

func (q *SimpleQueue) Enqueue(m Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.messages = append(q.messages, m)
	sort.SliceStable(q.messages, func(i int, j int) bool {
		return q.messages[i].Expires.Before(q.messages[j].Expires)
	})
	return nil
}

func (q *SimpleQueue) Dequeue() (Message, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.messages) == 0 {
		return Message{}, fmt.Errorf("queue is empty")
	}
	m := q.messages[0]
	q.messages = q.messages[1:]
	return m, nil
}

func (q *SimpleQueue) Peek() (Message, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if len(q.messages) == 0 {
		return Message{}, fmt.Errorf("queue is empty")
	}
	m := q.messages[0]
	return m, nil
}

func (q *SimpleQueue) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.messages) == 0
}
