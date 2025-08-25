package ws

import (
	"encoding/json"
	"errors"
	"pusher/internal/types"
	"pusher/pkg/logger"
	"sync"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type Subscriber struct {
	mu        sync.RWMutex    // 保护状态和连接
	isClosed  bool            // 原子性关闭状态
	conn      *websocket.Conn // 客户端连接
	sendChan  chan types.Data // 有缓冲的推送通道
	closeChan chan struct{}   // 关闭信号通道
}

// NewSubscriber 创建并初始化订阅者
func NewSubscriber(conn *websocket.Conn) *Subscriber {
	s := &Subscriber{
		conn:      conn,
		sendChan:  make(chan types.Data, types.SubscriberSendChanSize),
		closeChan: make(chan struct{}),
	}

	go s.writeLoop()

	return s
}

// Send 线程安全的消息发送
func (s *Subscriber) Send(data *types.Data) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.isClosed {
		return errors.New("subscriber closed")
	}

	select {
	case s.sendChan <- *data:
		return nil
	default:
		return errors.New("send buffer full")
	}
}

// writeLoop 处理发送消息的协程
func (s *Subscriber) writeLoop() {
	for {
		select {
		case data := <-s.sendChan:
			msg, err := json.Marshal(data)
			if err != nil {
				logger.GetLogger().Warnf("[subscriber] marshal error: %v", err)
				continue
			}
			err = s.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.GetLogger().Warnf("[subscriber] write message error: %v", err)
				continue
			}
		case <-s.closeChan:
			logger.GetLogger().Infof("[subscriber] closing subscriber for %s", s.conn.RemoteAddr())
			s.safeClose()
			return
		}
	}
}

// Close 关闭订阅者
func (s *Subscriber) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isClosed {
		return
	}

	close(s.closeChan)
	s.isClosed = true

	if s.conn != nil {
		_ = s.conn.Close()
	}
}

// GetCloseState 获取关闭状态
func (s *Subscriber) GetCloseState() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isClosed
}

// safeClose 关闭连接并清理资源
func (s *Subscriber) safeClose() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isClosed {
		s.isClosed = true
		if s.conn != nil {
			_ = s.conn.Close()
		}
	}
}
