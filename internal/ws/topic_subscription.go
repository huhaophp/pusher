package ws

import (
	"context"
	"pusher/internal/source"
	"pusher/internal/types"
	"pusher/pkg/logger"
	"sync"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type TopicSubscription struct {
	Topic       string
	subscribers map[string]*subscriberGroup
	mu          sync.RWMutex
	source      source.Source
	ctx         context.Context
	cancel      context.CancelFunc
}

type subscriberGroup struct {
	subs map[*websocket.Conn]*Subscriber
	mu   sync.RWMutex
}

func NewTopicSubscription(topic string, redisSource source.Source) *TopicSubscription {
	ctx, cancel := context.WithCancel(context.Background())
	ts := &TopicSubscription{
		Topic:       topic,
		subscribers: make(map[string]*subscriberGroup),
		source:      redisSource,
		ctx:         ctx,
		cancel:      cancel,
	}

	go ts.start()
	return ts
}

func (ts *TopicSubscription) start() {
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger().Errorf("topic %s panic: %v", ts.Topic, r)
		}
	}()

	err := ts.source.PullMessage(ts.ctx, ts.Topic, ts.onMessage)
	if err != nil {
		logger.GetLogger().Warnf("error pulling message: %v", err)
	}
}

func (ts *TopicSubscription) onMessage(data *types.Data) {
	ts.mu.RLock()
	group, exists := ts.subscribers[data.Type]
	ts.mu.RUnlock()

	if !exists {
		return
	}

	group.mu.RLock()
	defer group.mu.RUnlock()

	for conn, sub := range group.subs {
		select {
		case <-sub.closeChan:
			go ts.Remove(data.Type, conn) // 异步清理
		default:
			err := sub.Send(data)
			if err != nil {
				logger.GetLogger().Errorf("error sending message: %v", err)
			}
		}
	}
}

func (ts *TopicSubscription) Add(typ string, conn *websocket.Conn) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.subscribers[typ]; !exists {
		ts.subscribers[typ] = &subscriberGroup{
			subs: make(map[*websocket.Conn]*Subscriber),
		}
	}

	ts.subscribers[typ].mu.Lock()
	defer ts.subscribers[typ].mu.Unlock()
	ts.subscribers[typ].subs[conn] = NewSubscriber(conn)
}

func (ts *TopicSubscription) Remove(typ string, conn *websocket.Conn) {
	ts.mu.RLock()
	group, exists := ts.subscribers[typ]
	ts.mu.RUnlock()

	if !exists {
		return
	}

	group.mu.Lock()
	defer group.mu.Unlock()

	delete(group.subs, conn)

	if len(group.subs) == 0 {
		ts.mu.Lock()
		defer ts.mu.Unlock()
		delete(ts.subscribers, typ)
	}
}

func (ts *TopicSubscription) Close(conn *websocket.Conn) {
	ts.mu.RLock()
	subscribers := ts.subscribers
	ts.mu.RUnlock()
	for _, group := range subscribers {
		group.mu.RLock()
		if sub, ok := group.subs[conn]; ok {
			sub.Close()
			delete(group.subs, conn)
		}
		group.mu.RUnlock()
	}
}
