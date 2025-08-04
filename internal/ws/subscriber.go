package ws

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"pusher/internal/model"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type Subscriber struct {
	isClosed bool            // 是否关闭
	conn     *websocket.Conn // 客户端连接
	sendChan chan model.Data // 推送通道
}

// NewSubscriber creates a new subscriber.
func NewSubscriber(conn *websocket.Conn) *Subscriber {
	s := &Subscriber{
		isClosed: false,
		conn:     conn,
		sendChan: make(chan model.Data, 100),
	}
	go s.writeLoop()
	return s
}

// Send sends the message to the connection.
func (s *Subscriber) Send(data *model.Data) {
	select {
	case s.sendChan <- *data:
		log.Printf("[subscriber] sent message to conn %s", s.conn.RemoteAddr())
	default:
		log.Printf("[subscriber] send buffer full, dropping message for conn %s", s.conn.RemoteAddr())
	}
}

// writeLoop writes the message to the connection.
func (s *Subscriber) writeLoop() {
	for data := range s.sendChan {
		msg, err := json.Marshal(data)
		if err != nil {
			log.Printf("[subscriber] Marshal error: %v", err)
			continue
		}
		err = s.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("[subscriber] write failed: %v", err)
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
