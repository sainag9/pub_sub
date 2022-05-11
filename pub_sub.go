package pub_sub

import "sync"

type GoPS interface {
	Publish(topic, message string)
	Subscribe(topic string) chan string
	Unsubscribe(ch chan string)
}

type pubSub struct {
	lock     sync.RWMutex
	pubStore map[string][]chan string
	subStore map[chan string]string
	capacity int
}

func NewPubSub(capacity int) GoPS {
	pub := map[string][]chan string{}
	sub := map[chan string]string{}
	return &pubSub{pubStore: pub, subStore: sub, capacity: capacity}
}

func (ps *pubSub) Publish(topic, message string) {

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

func (ps *pubSub) Subscribe(topic string) chan string {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	ch := make(chan string, ps.capacity)
	if len(ps.pubStore[topic]) == 0 {
		ps.pubStore[topic] = make([]chan string, 0)
	}
	// add the channel to the list of subscribers for the corresponding topic
	ps.pubStore[topic] = append(ps.pubStore[topic], ch)
	// 1:1 mapping for channel and topic
	ps.subStore[ch] = topic
	return ch
}

func (ps *pubSub) Unsubscribe(ch chan string) {

	ps.lock.RLock()
	defer ps.lock.RUnlock()
	topic := ps.subStore[ch]
	updatedSubs := make([]chan string, 0)
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
