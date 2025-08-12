package ws

import (
	"encoding/json"
	"pusher/internal/types"
	"pusher/pkg/logger"
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type DefaultHandler struct {
	SubscriptionManager *SubscriptionManager
}

// OnOpen 当连接打开时调用
func (h *DefaultHandler) OnOpen(c *websocket.Conn) {
	logger.GetLogger().Debugf("connection opened, remote address: %s", c.RemoteAddr().String())
}

// OnMessage 当客户端发送消息时调用
func (h *DefaultHandler) OnMessage(c *websocket.Conn, _ websocket.MessageType, data []byte) {
	var request types.Request
	if err := json.Unmarshal(data, &request); err != nil {
		logger.GetLogger().Errorf("unmarshal failed, err: %+v, remote address: %s", err, c.RemoteAddr().String())
		return
	}

	logger.GetLogger().Debugf("received message: %+v, remote address: %s", request, c.RemoteAddr().String())

	handlers := map[string]func(*websocket.Conn, *types.Request){
		types.ActionSubscribe:   h.onSubscribe,
		types.ActionUnsubscribe: h.onUnsubscribe,
		types.ActionPing:        h.onPing,
	}

	handler, ok := handlers[request.Action]
	if !ok {
		logger.GetLogger().Errorf("invalid action: %s, remote address: %s", request.Action, c.RemoteAddr().String())
		return
	}

	handler(c, &request)
}

// OnClose 当连接关闭时调用
func (h *DefaultHandler) OnClose(c *websocket.Conn, err error) {
	h.SubscriptionManager.Close(c)
	logger.GetLogger().Debugf("connection closed, remote address: %s, err: %v", c.RemoteAddr().String(), err)
}

// onSubscribe 当客户端订阅主题时调用
func (h *DefaultHandler) onSubscribe(c *websocket.Conn, request *types.Request) {
	err := h.SubscriptionManager.Subscribe(request.Params.Topic, request.Params.Type, c)
	if err != nil {
		h.respondError(c, request, err)
		return
	}
	h.respondSuccess(c, request)
}

// onUnsubscribe 当客户端取消订阅主题时调用
func (h *DefaultHandler) onUnsubscribe(c *websocket.Conn, request *types.Request) {
	h.SubscriptionManager.Unsubscribe(request.Params.Topic, request.Params.Type, c)
	h.respondSuccess(c, request)
}

// onPing 当客户端发送ping消息时调用
func (h *DefaultHandler) onPing(c *websocket.Conn, request *types.Request) {
	h.respondSuccess(c, request)
}

// Response 发送响应到客户端
func (h *DefaultHandler) response(c *websocket.Conn, response *types.Response) {
	data, err := json.Marshal(response)
	if err != nil {
		logger.GetLogger().Errorf("json marshal failed, err: %+v", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logger.GetLogger().Errorf("write message failed, err: %+v", err)
	}
}

func (h *DefaultHandler) respondSuccess(c *websocket.Conn, request *types.Request) {
	h.response(c, &types.Response{
		Action:    request.Action,
		RequestID: request.RequestID,
		Timestamp: time.Now().UnixMilli(),
		Status:    types.StatusOK,
	})
}

func (h *DefaultHandler) respondError(c *websocket.Conn, request *types.Request, err error) {
	logger.GetLogger().Errorf("action %s failed, err: %+v", request.Action, err)
	h.response(c, &types.Response{
		Action:    request.Action,
		RequestID: request.RequestID,
		Timestamp: time.Now().UnixMilli(),
		Status:    types.StatusFail,
		ErrorMsg:  err.Error(),
	})
}
