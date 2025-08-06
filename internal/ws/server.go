package ws

import (
	"fmt"
	"net/http"
	"pusher/config"
	"pusher/pkg/logger"
	"sync"
	"time"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

// Handler 事件处理接口
type Handler interface {
	OnOpen(c *websocket.Conn)
	OnMessage(c *websocket.Conn, messageType websocket.MessageType, data []byte)
	OnClose(c *websocket.Conn, err error)
}

// WebsocketServer WebSocket服务器
type WebsocketServer struct {
	upgrader *websocket.Upgrader
	config   *config.APP
	handler  Handler
	conns    map[*websocket.Conn]struct{}
	mu       sync.RWMutex
}

// NewWebsocketServer 创建一个新的WebSocket服务器实例
func NewWebsocketServer(config *config.APP, handler Handler) *WebsocketServer {
	if handler == nil {
		handler = &DefaultHandler{}
	}

	upgrader := websocket.NewUpgrader()

	ws := &WebsocketServer{
		config:   config,
		handler:  handler,
		upgrader: upgrader,
		conns:    make(map[*websocket.Conn]struct{}),
	}

	go ws.monitor()

	// 注册回调
	ws.upgrader.OnOpen(ws.onOpen)
	ws.upgrader.OnClose(ws.onClose)
	ws.upgrader.OnMessage(ws.onMessage)

	return ws
}

// onWebsocket 升级连接到WebSocket
func (ws *WebsocketServer) onWebsocket(w http.ResponseWriter, r *http.Request) {
	if _, err := ws.upgrader.Upgrade(w, r, nil); err != nil {
		logger.Infof("upgrade failed: %v", err)
		http.Error(w, "upgrade failed", http.StatusInternalServerError)
	}
}

// onOpen 当连接打开时调用
func (ws *WebsocketServer) onOpen(c *websocket.Conn) {
	ws.addConn(c)
	ws.handler.OnOpen(c)
}

// onMessage 当客户端发送消息时调用
func (ws *WebsocketServer) onMessage(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
	ws.handler.OnMessage(c, messageType, data)
}

// onClose 当连接关闭时调用
func (ws *WebsocketServer) onClose(c *websocket.Conn, err error) {
	ws.remConn(c)
	ws.handler.OnClose(c, err)
}

// addConn 添加连接
func (ws *WebsocketServer) addConn(conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.conns[conn] = struct{}{}
}

// remConn 删除连接
func (ws *WebsocketServer) remConn(conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	delete(ws.conns, conn)
}

// GetConnAll 获取所有连接
func (ws *WebsocketServer) GetConnAll() map[*websocket.Conn]struct{} {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.conns
}

// monitor 监控连接
func (ws *WebsocketServer) monitor() {
	for {
		ws.mu.RLock()
		logger.Infof("current connections: %d, connections: %+v", len(ws.conns), ws.conns)
		ws.mu.RUnlock()
		time.Sleep(time.Second * 5)
	}
}

// Run 启动WebSocket服务器
func (ws *WebsocketServer) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", ws.onWebsocket)

	engine := nbhttp.NewEngine(nbhttp.Config{
		Name:                    ws.config.Name,
		Network:                 "tcp",
		Addrs:                   []string{fmt.Sprintf(":%s", ws.config.Port)},
		MaxLoad:                 1000000,
		ReleaseWebsocketPayload: true,
		Handler:                 mux,
	})

	if err := engine.Start(); err != nil {
		return err
	}

	logger.Infof("ws server started at %s", ws.config.Port)

	return nil
}
