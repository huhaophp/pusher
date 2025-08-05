package ws

import (
	"pusher/internal/source"
	"pusher/pkg/logger"
	"sync"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

// SubscriptionManager manages the topics and their subscribers.
type SubscriptionManager struct {
	topics      map[string]*TopicSubscription
	redisSource *source.RedisSource
	mu          sync.RWMutex
}

// NewSubscriptionManager creates a new subscription manager.
// It initializes the topics and starts the monitor.
func NewSubscriptionManager(topics []string, redisSource *source.RedisSource) *SubscriptionManager {
	subscriptionManager := &SubscriptionManager{
		topics: make(map[string]*TopicSubscription),
	}

	for _, topic := range topics {
		subscriptionManager.topics[topic] = NewTopicSubscription(topic, redisSource)
	}

	go subscriptionManager.monitor()

	return subscriptionManager
}

// Subscribe adds a new subscriber to the topic.
func (sm *SubscriptionManager) Subscribe(topic, typ string, c *websocket.Conn) {
	sm.mu.Lock()
	sm.topics[topic].Add(typ, c)
	sm.mu.Unlock()
}

// Unsubscribe removes a subscriber from the topic.
func (sm *SubscriptionManager) Unsubscribe(topic, typ string, conn *websocket.Conn) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if t, ok := sm.topics[topic]; ok {
		t.Remove(typ, conn)
	}
}

// monitor monitors the topics and their subscribers.
func (sm *SubscriptionManager) monitor() {
	for {
		sm.mu.RLock()
		for topic, t := range sm.topics {
			for typ, subscribers := range t.subscribers {
				logger.Infof("topic: %s, type: %s, subscribers: %d", topic, typ, len(subscribers))
			}
		}
		sm.mu.RUnlock()
		time.Sleep(time.Second * 5)
	}
}
