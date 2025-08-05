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

func (h *DefaultHandler) OnOpen(c *websocket.Conn) {
	logger.Infof("connection opened: %s", c.RemoteAddr().String())
}

func (h *DefaultHandler) OnMessage(c *websocket.Conn, _ websocket.MessageType, data []byte) {
	var request types.Request
	if err := json.Unmarshal(data, &request); err != nil {
		logger.Errorf("unmarshal failed, err: %+v", err)
		return
	}

	if !slices.Contains(types.AllAction, request.Action) {
		logger.Errorf("invalid action: %s", request.Action)
		return
	}

	handlers := map[string]func(*websocket.Conn, *types.Request){
		types.ActionSubscribe:   h.onSubscribe,
		types.ActionUnsubscribe: h.onUnsubscribe,
		types.ActionPing:        h.onPing,
	}

	handlers[request.Action](c, &request)
}

func (h *DefaultHandler) OnClose(c *websocket.Conn, err error) {
	logger.Infof("connection closed: %s, err: %v", c.RemoteAddr().String(), err)
}

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
