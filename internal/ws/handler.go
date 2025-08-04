package ws

import (
	"encoding/json"
	"log"
	"pusher/internal/model"
	"slices"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type DefaultHandler struct {
	SubscriptionManager *SubscriptionManager
}

func (h *DefaultHandler) OnOpen(c *websocket.Conn) {
	log.Printf("连接打开: %s", c.RemoteAddr().String())
}

func (h *DefaultHandler) OnMessage(c *websocket.Conn, _ websocket.MessageType, data []byte) {
	var request model.Request
	if err := json.Unmarshal(data, &request); err != nil {
		log.Printf("unmarshal failed, err: %+v", err)
		return
	}

	if !slices.Contains(model.AllAction, request.Action) {
		log.Printf("invalid action: %s", request.Action)
		return
	}

	handlers := map[string]func(*websocket.Conn, *model.Request){
		model.ActionSubscribe:   h.onSubscribe,
		model.ActionUnsubscribe: h.onUnsubscribe,
		model.ActionPing:        h.onPing,
	}

	handlers[request.Action](c, &request)
}

func (h *DefaultHandler) OnClose(c *websocket.Conn, err error) {
	log.Printf("连接关闭: %s, 错误: %v", c.RemoteAddr().String(), err)
}

func (h *DefaultHandler) onSubscribe(c *websocket.Conn, request *model.Request) {
	h.SubscriptionManager.Subscribe(request.Params.Topic, request.Params.Type, c)

	resp := model.Response{
		Action:    request.Action,
		RequestID: request.RequestID,
		Timestamp: time.Now().UnixMilli(),
		Status:    model.StatusOK,
	}

	h.Response(c, &resp)
}

func (h *DefaultHandler) onUnsubscribe(c *websocket.Conn, request *model.Request) {
	h.SubscriptionManager.Unsubscribe(request.Params.Topic, request.Params.Type, c)

	resp := model.Response{
		Action:    request.Action,
		RequestID: request.RequestID,
		Timestamp: time.Now().UnixMilli(),
		Status:    model.StatusOK,
	}

	h.Response(c, &resp)
}

func (h *DefaultHandler) onPing(c *websocket.Conn, request *model.Request) {
	c.SetDeadline(time.Now().Add(time.Minute * 2))
	resp := model.Response{
		Action:    request.Action,
		RequestID: request.RequestID,
		Timestamp: time.Now().UnixMilli(),
		Status:    model.StatusOK,
	}
	h.Response(c, &resp)
}

func (h *DefaultHandler) Response(c *websocket.Conn, response *model.Response) {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("marshal failed, err: %+v", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("write message failed, err: %+v", err)
	}
}
