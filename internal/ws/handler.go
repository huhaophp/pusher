package ws

import (
	"encoding/json"
	"pusher/internal/types"
	"pusher/pkg/logger"
	"slices"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type DefaultHandler struct {
	SubscriptionManager *SubscriptionManager
}

// OnOpen 当连接打开时调用
func (h *DefaultHandler) OnOpen(c *websocket.Conn) {
	logger.Debugf("connection opened, remote address: %s", c.RemoteAddr().String())
}

// OnMessage 当客户端发送消息时调用
func (h *DefaultHandler) OnMessage(c *websocket.Conn, _ websocket.MessageType, data []byte) {
	var request types.Request
	if err := json.Unmarshal(data, &request); err != nil {
		logger.Errorf("unmarshal failed, err: %+v, remote address: %s", err, c.RemoteAddr().String())
		return
	}

	logger.Debugf("received message: %+v, remote address: %s", request, c.RemoteAddr().String())

	if !slices.Contains(types.AllAction, request.Action) {
		logger.Errorf("invalid action: %s, remote address: %s", request.Action, c.RemoteAddr().String())
		return
	}

	handlers := map[string]func(*websocket.Conn, *types.Request){
		types.ActionSubscribe:   h.onSubscribe,
		types.ActionUnsubscribe: h.onUnsubscribe,
		types.ActionPing:        h.onPing,
	}

	handlers[request.Action](c, &request)
}

// OnClose 当连接关闭时调用
func (h *DefaultHandler) OnClose(c *websocket.Conn, err error) {
	logger.Debugf("connection closed, remote address: %s, err: %v", c.RemoteAddr().String(), err)
}

// onSubscribe 当客户端订阅主题时调用
func (h *DefaultHandler) onSubscribe(c *websocket.Conn, request *types.Request) {
	h.SubscriptionManager.Subscribe(request.Params.Topic, request.Params.Type, c)

	resp := types.Response{
		Action:    request.Action,
		RequestID: request.RequestID,
		Timestamp: time.Now().UnixMilli(),
		Status:    types.StatusOK,
	}

	h.Response(c, &resp)
}

// onUnsubscribe 当客户端取消订阅主题时调用
func (h *DefaultHandler) onUnsubscribe(c *websocket.Conn, request *types.Request) {
	h.SubscriptionManager.Unsubscribe(request.Params.Topic, request.Params.Type, c)

	resp := types.Response{
		Action:    request.Action,
		RequestID: request.RequestID,
		Timestamp: time.Now().UnixMilli(),
		Status:    types.StatusOK,
	}

	h.Response(c, &resp)
}

// onPing 当客户端发送ping消息时调用
func (h *DefaultHandler) onPing(c *websocket.Conn, request *types.Request) {
	err := c.SetDeadline(time.Now().Add(types.ConnDeadlineTime))
	if err != nil {
		logger.Errorf("set deadline failed, err: %+v", err)
		h.Response(c, &types.Response{
			Action:    request.Action,
			RequestID: request.RequestID,
			Timestamp: time.Now().UnixMilli(),
			Status:    types.StatusFail,
			ErrorMsg:  err.Error(),
		})
		return
	}
	resp := types.Response{
		Action:    request.Action,
		RequestID: request.RequestID,
		Timestamp: time.Now().UnixMilli(),
		Status:    types.StatusOK,
	}
	h.Response(c, &resp)
}

// Response 发送响应到客户端
func (h *DefaultHandler) Response(c *websocket.Conn, response *types.Response) {
	data, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("json marshal failed, err: %+v", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logger.Errorf("write message failed, err: %+v", err)
	}
}
