package pub_sub

import (
	"sync"
)

type GoPS[T any] interface {
	Publish(topic string, message T)
	Subscribe(topic string) chan T
	Unsubscribe(ch chan T)
}

type pubSub[T any] struct {
	lock     sync.RWMutex
	pubStore map[string][]chan T
	subStore map[chan T]string
	capacity int
}

func NewPubSub[T any](capacity int) GoPS[T] {
	pub := map[string][]chan T{}
	sub := map[chan T]string{}
	return &pubSub[T]{pubStore: pub, subStore: sub, capacity: capacity}
}

func (ps *pubSub[T]) Publish(topic string, message T) {

	ps.lock.RLock()
	defer ps.lock.RUnlock()

	// if not found return
	if _, ok := ps.pubStore[topic]; !ok {
		return
	}
	for _, ch := range ps.pubStore[topic] {
		ch <- message
	}
}

func (ps *pubSub[T]) Subscribe(topic string) chan T {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	ch := make(chan T, ps.capacity)
	if len(ps.pubStore[topic]) == 0 {
		ps.pubStore[topic] = make([]chan T, 0)
	}
	// add the channel to the list of subscribers for the corresponding topic
	ps.pubStore[topic] = append(ps.pubStore[topic], ch)
	// 1:1 mapping for channel and topic
	ps.subStore[ch] = topic
	return ch
}

func (ps *pubSub[T]) Unsubscribe(ch chan T) {

	ps.lock.RLock()
	defer ps.lock.RUnlock()
	topic := ps.subStore[ch]
	updatedSubs := make([]chan T, 0)
	for _, oldCh := range ps.pubStore[topic] {
		if oldCh != ch {
			updatedSubs = append(updatedSubs, oldCh)
		}
	}
	// updatedSubs doesn't contain the give channel ch
	ps.pubStore[topic] = updatedSubs
	// delete 1:1 mapping for the channel, and it's topic
	delete(ps.subStore, ch)
}
