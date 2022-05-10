package pub_sub

type GoPS interface {
	Publish(key, value string)
	Subscribe(key string) chan string
	Unsubscribe(key string)
}

type pubSub struct {
	store map[string]chan string
}

const (
	maxChannelSize = 5
)

func NewPubSub() GoPS {
	m := map[string]chan string{}
	return &pubSub{store: m}
}

func (ps *pubSub) Publish(key, value string) {

	// if not found return
	if _, ok := ps.store[key]; !ok {
		return
	}
	ps.store[key] <- value
}

func (ps *pubSub) Subscribe(key string) chan string {
	ch := make(chan string, maxChannelSize)
	ps.store[key] = ch
	return ch
}

func (ps *pubSub) Unsubscribe(key string) {
	ps.store[key] = nil
}
