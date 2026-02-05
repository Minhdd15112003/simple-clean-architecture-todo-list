package pubsub

import (
	"context"
	"log"
	"social-todo-list/middleware"
	"sync"
)

type localPubSub struct {
	messageQueue chan *Message
	mapChannel   map[Topic][]chan *Message
	locker       *sync.RWMutex
}

func NewPubSub() *localPubSub {
	pb := &localPubSub{
		messageQueue: make(chan *Message, 10000),
		mapChannel:   make(map[Topic][]chan *Message),
		locker:       new(sync.RWMutex),
	}
	pb.run()

	return pb
}

func (ps *localPubSub) Publish(ctx context.Context, topic Topic, data *Message) error {
	data.SetChannel(topic)

	go func() {
		defer middleware.RecoverGoroutine()
		ps.messageQueue <- data
		log.Println("New message published", data.String())
	}()
	return nil
}

func (ps *localPubSub) Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, unsubscribe func()) {
	c := make(chan *Message)

	ps.locker.Lock()

	if val, ok := ps.mapChannel[topic]; ok {
		val = append(val, c)
		ps.mapChannel[topic] = val
	} else {
		ps.mapChannel[topic] = []chan *Message{c}
	}
	ps.locker.Unlock()

	return c, func() {
		log.Println("Unsubscribe")

		if chans, ok := ps.mapChannel[topic]; ok {
			for i := range chans {
				if chans[i] == c {
					chans = append(chans[:i], chans[i+1:]...)

					ps.locker.Lock()
					ps.mapChannel[topic] = chans
					ps.locker.Unlock()

					close(c)
					break
				}
			}
		}
	}
}

func (ps *localPubSub) run() error {
	go func() {
		defer middleware.RecoverGoroutine()

		for {
			mess := <-ps.messageQueue
			log.Println("Message dequeue", mess.String())

			ps.locker.RLock()

			if subs, ok := ps.mapChannel[mess.Channel()]; ok {
				for i := range subs {
					go func(c chan *Message) {
						defer middleware.RecoverGoroutine()
						c <- mess
					}(subs[i])
				}
			}
			ps.locker.RUnlock()
		}
	}()
	return nil
}
