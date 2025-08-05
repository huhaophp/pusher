package ws

import (
	"encoding/json"
	"errors"
	"net"
	"pusher/internal/types"
	"pusher/pkg/logger"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type Subscriber struct {
	isClosed bool            // 是否关闭
	conn     *websocket.Conn // 客户端连接
	sendChan chan types.Data // 推送通道
}

// NewSubscriber creates a new subscriber.
func NewSubscriber(conn *websocket.Conn) *Subscriber {
	s := &Subscriber{
		isClosed: false,
		conn:     conn,
		sendChan: make(chan types.Data, types.SubscriberSendChanSize),
	}
	go s.writeLoop()
	return s
}

// Send sends the message to the connection.
func (s *Subscriber) Send(data *types.Data) {
	select {
	case s.sendChan <- *data:
	default:
		logger.Infof("[subscriber] send buffer full, dropping message for conn %s", s.conn.RemoteAddr())
	}
}

// writeLoop writes the message to the connection.
func (s *Subscriber) writeLoop() {
	for data := range s.sendChan {
		msg, err := json.Marshal(data)
		if err != nil {
			logger.Infof("[subscriber] Marshal error: %v", err)
			continue
		}
		err = s.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			logger.Infof("[subscriber] write failed: %v", err)
			if errors.Is(err, net.ErrClosed) {
				s.isClosed = true
				return
			}
		}
	}
}

// GetCloseState gets the close state of the connection.
func (s *Subscriber) GetCloseState() bool {
	return s.isClosed
}
