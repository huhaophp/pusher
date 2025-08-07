package ws

import (
	"fmt"
	"pusher/internal/source"
	"pusher/pkg/logger"
	"sync"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

// SubscriptionManager 管理主题和订阅者
type SubscriptionManager struct {
	subscriptions map[string]*TopicSubscription
	mu            sync.RWMutex
}

// NewSubscriptionManager 创建一个新的订阅管理器
func NewSubscriptionManager(redisTopicPuller *source.TopicPuller, kafkaTopicPuller *source.TopicPuller) *SubscriptionManager {
	subscriptionManager := &SubscriptionManager{
		subscriptions: make(map[string]*TopicSubscription),
	}

	for _, topic := range redisTopicPuller.Topics {
		subscriptionManager.subscriptions[topic] = NewTopicSubscription(topic, redisTopicPuller.Source)
	}

	for _, topic := range kafkaTopicPuller.Topics {
		subscriptionManager.subscriptions[topic] = NewTopicSubscription(topic, kafkaTopicPuller.Source)
	}

	go subscriptionManager.monitor()

	return subscriptionManager
}

// Subscribe 添加一个新的订阅者到主题
func (sm *SubscriptionManager) Subscribe(topic, typ string, c *websocket.Conn) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sub, ok := sm.subscriptions[topic]
	if !ok {
		logger.Warnf("subscribe called for unknown topic: %s", topic)
		return fmt.Errorf("subscribe called for unknown topic: %s", topic)
	}
	sub.Add(typ, c)
	return nil
}

// Unsubscribe 从主题中删除一个订阅者
func (sm *SubscriptionManager) Unsubscribe(topic, typ string, conn *websocket.Conn) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if t, ok := sm.subscriptions[topic]; ok {
		t.Remove(typ, conn)
	}
}

// monitor 监控主题和订阅者
func (sm *SubscriptionManager) monitor() {
	for {
		sm.mu.RLock()
		for topic, t := range sm.subscriptions {
			for typ, subscribers := range t.subscribers {
				logger.Infof("topic: %s, type: %s, subscribers: %d", topic, typ, len(subscribers))
			}
		}
		sm.mu.RUnlock()
		time.Sleep(time.Second * 5)
	}
}
