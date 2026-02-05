package pubsub

import (
	"context"
	"fmt"
	"time"
)

type Topic string

type PubSup interface {
	Publish(ctx context.Context, topic Topic, data *Message) error
	Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, close func())
}

type Message struct {
	id        string
	topic     Topic
	data      interface{}
	CreatedAt time.Time
}

func NewMessage(data interface{}) *Message {
	now := time.Now().UTC()

	return &Message{
		id:        fmt.Sprintf("%d", now.UnixNano()),
		data:      data,
		CreatedAt: now,
	}
}

func (evt *Message) String() string {
	return fmt.Sprintf("Message %s value %v", evt.topic, evt.data)
}

func (evt *Message) Channel() Topic {
	return evt.topic
}

func (evt *Message) SetChannel(topic Topic) {
	evt.topic = topic
}

func (evt *Message) Data() interface{} {
	return evt.data
}
