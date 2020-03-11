package pubsub

import (
	"strconv"
	"sync"
	"time"
)

func NewPublisher() Publisher {
	return &publisher{
		queue:         NewSimpleQueue(),
		nextMessageID: GeneratePrefixedID("msg-"),
		nextSubID:     GeneratePrefixedID("sub-"),
		subscriptions: map[string]*subscription{},
	}
}

type publisher struct {
	queue         Queue
	nextMessageID func() string
	nextSubID     func() string

	mu            sync.RWMutex
	subscriptions map[string]*subscription
}

func (p *publisher) Publish(topic Topic, data interface{}) error {
	return p.publish(Message{
		ID:      p.nextMessageID(),
		Topic:   topic,
		Data:    data,
		Expires: time.Now(),
	})
}

func (p *publisher) publish(m Message) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, s := range p.subscriptions {
		if s.filter.Accept(m.Topic) {
			s.callback(m.Topic, m.Data)
		}
	}
	return nil
}

func (p *publisher) Subscribe(topic Topic, callback func(topic Topic, data interface{})) Subscription {
	p.mu.Lock()
	defer p.mu.Unlock()
	id := p.nextSubID()
	s := &subscription{publisher: p, id: id, filter: topic, callback: callback}
	p.subscriptions[id] = s
	return s
}

func (p *publisher) cancelSubscription(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.subscriptions, id)
}

type subscription struct {
	publisher *publisher
	id        string
	filter    Topic
	callback  func(topic Topic, data interface{})
}

func (s *subscription) Cancel() {
	s.publisher.cancelSubscription(s.id)
}

type Publisher interface {
	Publish(topic Topic, data interface{}) error
	Subscribe(topic Topic, callback func(topic Topic, data interface{})) Subscription
}

type Subscription interface {
	Cancel()
}

type Message struct {
	ID      string
	Topic   Topic
	Data    interface{}
	Expires time.Time
	Retain  bool
}

func GeneratePrefixedID(prefix string) func() string {
	pre := []byte(prefix + ":")
	var c int64
	var mu sync.Mutex
	return func() string {
		mu.Lock()
		defer mu.Unlock()
		c++
		return string(strconv.AppendInt(pre, c, 10))
	}
}
