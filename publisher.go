package pubsub

import (
	"sync"
)

func NewPublisher() Publisher {
	return &publisher{
		subscriptions: map[int]*subscription{},
	}
}

type publisher struct {
	mu            sync.RWMutex
	next          int
	subscriptions map[int]*subscription
}

func (p *publisher) Publish(topic Topic, args ...interface{}) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, s := range p.subscriptions {
		if s.filter.Accept(topic) {
			s.callback(topic, args...)
		}
	}
}

func (p *publisher) Subscribe(topic Topic, callback func(Topic, ...interface{})) Subscription {
	p.mu.Lock()
	defer p.mu.Unlock()
	id := p.next
	s := &subscription{publisher: p, id: id, filter: topic, callback: callback}
	p.subscriptions[id] = s
	p.next++
	return s
}

func (p *publisher) cancelSubscription(id int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.subscriptions, id)
}

type subscription struct {
	publisher *publisher
	id        int
	filter    Topic
	callback  func(Topic, ...interface{})
}

func (s *subscription) Cancel() {
	s.publisher.cancelSubscription(s.id)
}

type Publisher interface {
	Publish(topic Topic, args ...interface{})
	Subscribe(topic Topic, callback func(Topic, ...interface{})) Subscription
}

type Subscription interface {
	Cancel()
}
