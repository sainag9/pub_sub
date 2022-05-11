# pub_sub

## usage

```golang

	client := NewPubSub(10)
	ch1 := client.Subscribe("foo")
	ch2 := client.Subscribe("foo")
	ch3 := client.Subscribe("foo")
	client.Publish("foo", "hello foofers!")
	
	fmt.Println(<-ch1, <-ch2, <-ch3)
```