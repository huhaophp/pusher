package ws

import (
	"context"
	"pusher/internal/model"
	"pusher/internal/source"
	"sync"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type TopicSubscription struct {
	Topic       string
	subscribers map[string]map[*websocket.Conn]*Subscriber
	mu          sync.RWMutex
	source      *source.RedisSource
	ch          chan []byte
}

// NewTopicSubscription creates a new topic subscription.
func NewTopicSubscription(topic string, redisSource *source.RedisSource) *TopicSubscription {
	topicSubscription := &TopicSubscription{
		Topic:       topic,
		source:      redisSource,
		subscribers: make(map[string]map[*websocket.Conn]*Subscriber),
	}

	go topicSubscription.start()

	return topicSubscription
}

// start starts the topic subscription.
func (ts *TopicSubscription) start() {
	ts.source.PullMessage(context.Background(), ts.Topic, ts.onMessage)
}

// onMessage handles the message from the source.
func (ts *TopicSubscription) onMessage(data *model.Data) {
	for conn, subscriber := range ts.subscribers[data.Type] {
		if subscriber.isClosed {
			ts.Remove(data.Type, conn)
			continue
		}
		subscriber.Send(data)
	}
}

// Add adds a new subscriber to the topic.
func (ts *TopicSubscription) Add(typ string, c *websocket.Conn) {
	ts.mu.Lock()
	if _, ok := ts.subscribers[typ]; !ok {
		ts.subscribers[typ] = make(map[*websocket.Conn]*Subscriber)
	}
	ts.subscribers[typ][c] = NewSubscriber(c)
	ts.mu.Unlock()
}

// Remove removes a subscriber from the topic.
func (ts *TopicSubscription) Remove(typ string, conn *websocket.Conn) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.subscribers[typ], conn)
}
