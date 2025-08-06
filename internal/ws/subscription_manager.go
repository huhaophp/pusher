package ws

import (
	"pusher/internal/source"
	"pusher/pkg/logger"
	"sync"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

// SubscriptionManager 管理主题和订阅者
type SubscriptionManager struct {
	topics map[string]*TopicSubscription
	mu     sync.RWMutex
}

// NewSubscriptionManager 创建一个新的订阅管理器
// 它初始化主题并启动监控
func NewSubscriptionManager(topics []string, source source.Source) *SubscriptionManager {
	subscriptionManager := &SubscriptionManager{
		topics: make(map[string]*TopicSubscription),
	}

	for _, topic := range topics {
		subscriptionManager.topics[topic] = NewTopicSubscription(topic, source)
	}

	go subscriptionManager.monitor()

	return subscriptionManager
}

// Subscribe 添加一个新的订阅者到主题
func (sm *SubscriptionManager) Subscribe(topic, typ string, c *websocket.Conn) {
	sm.mu.Lock()
	sm.topics[topic].Add(typ, c)
	sm.mu.Unlock()
}

// Unsubscribe 从主题中删除一个订阅者
func (sm *SubscriptionManager) Unsubscribe(topic, typ string, conn *websocket.Conn) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if t, ok := sm.topics[topic]; ok {
		t.Remove(typ, conn)
	}
}

// monitor 监控主题和订阅者
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
