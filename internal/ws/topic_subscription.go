package ws

import (
	"context"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"pusher/internal/source"
	"pusher/internal/types"
	"sync"
)

type TopicSubscription struct {
	Topic       string
	subscribers map[string]map[*websocket.Conn]*Subscriber
	mu          sync.RWMutex
	source      source.Source
}

// NewTopicSubscription 创建一个新的主题订阅
func NewTopicSubscription(topic string, redisSource source.Source) *TopicSubscription {
	topicSubscription := &TopicSubscription{
		Topic:       topic,
		source:      redisSource,
		subscribers: make(map[string]map[*websocket.Conn]*Subscriber),
	}

	go topicSubscription.start()

	return topicSubscription
}

// start 启动主题订阅
func (ts *TopicSubscription) start() {
	ts.source.PullMessage(context.Background(), ts.Topic, ts.onMessage)
}

// onMessage 处理来自源的消息
func (ts *TopicSubscription) onMessage(data *types.Data) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	for conn, subscriber := range ts.subscribers[data.Type] {
		if subscriber.isClosed {
			ts.Remove(data.Type, conn)
			continue
		}
		subscriber.Send(data)
	}
}

// Add 添加一个新的订阅者到主题
func (ts *TopicSubscription) Add(typ string, c *websocket.Conn) {
	ts.mu.Lock()
	if _, ok := ts.subscribers[typ]; !ok {
		ts.subscribers[typ] = make(map[*websocket.Conn]*Subscriber)
	}
	ts.subscribers[typ][c] = NewSubscriber(c)
	ts.mu.Unlock()
}

// Remove 从主题中删除一个订阅者
func (ts *TopicSubscription) Remove(typ string, conn *websocket.Conn) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.subscribers[typ], conn)
}
