package pub_sub

import (
	"testing"
)

func Test_single_subscriber(t *testing.T) {

	pubMsg := "hello foo!"
	client := NewPubSub(10)

	ch := client.Subscribe("foo")
	// 1 message
	client.Publish("foo", pubMsg)

	if v := <-ch; v != pubMsg {
		t.Errorf("expected %s, got %s", pubMsg, v)
	}
	// send multiple messages

	msgs := []string{"Hello", "ola", "ciao"}

	for _, v := range msgs {
		client.Publish("foo", v)
	}

	// order is maintained
	for _, msg := range msgs {

		if v := <-ch; v != msg {
			t.Errorf("expected %s, got %s", msg, v)
		}
	}
}

func Test_multiple_subscriber(t *testing.T) {

	pubMsg := "hello foofers!"
	client := NewPubSub(10)
	ch1 := client.Subscribe("foo")
	ch2 := client.Subscribe("foo")
	ch3 := client.Subscribe("foo")

	client.Publish("foo", pubMsg)

	if <-ch1 != pubMsg || <-ch2 != pubMsg || <-ch3 != pubMsg {
		t.Error("received messages not matching")
	}
}

func Test_unsubscribe(t *testing.T) {

	pubMsg := "hello foof!"
	client := NewPubSub(10)
	ch1 := client.Subscribe("foo")

	client.Publish("foo", pubMsg)

	if v := <-ch1; v != pubMsg {
		t.Errorf("got %s expected %s", v, pubMsg)
	}
	client.Unsubscribe(ch1)

}

func Benchmark_basic_func(b *testing.B) {
	client := NewPubSub(1000000)
	ch := client.Subscribe("foo")
	pubMsg := "hello foo!"
	for i := 0; i < b.N; i++ {
		client.Publish("foo", pubMsg)
		if v, ok := <-ch; v != pubMsg && !ok {
			b.Error("message mismatching")
		}
	}
	close(ch)
}
