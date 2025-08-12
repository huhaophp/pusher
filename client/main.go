package main

import (
	"encoding/json"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Subscribe struct {
	RequestId string `json:"request_id"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
	Params    struct {
		Topic string `json:"topic"`
		Type  string `json:"type"`
	} `json:"params"`
}

func newClient(id int, wsURL string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Printf("[Client %d] ❌ 连接失败: %v", id, err)
		return
	}
	defer conn.Close()
	log.Printf("[Client %d] ✅ 已连接", id)

	// 发送订阅
	sub := Subscribe{
		RequestId: time.Now().Format("150405"),
		Action:    "subscribe",
		Timestamp: time.Now().UnixMilli(),
	}
	sub.Params.Topic = "topic1"
	sub.Params.Type = "test"

	subBytes, _ := json.Marshal(sub)
	if err := conn.WriteMessage(websocket.TextMessage, subBytes); err != nil {
		log.Printf("[Client %d] ❌ 发送订阅失败: %v", id, err)
		return
	}

	// ping 定时器
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	// 读协程
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[Client %d] ❌ 读取错误: %v", id, err)
				return
			}
			//log.Printf("[Client %d] ✅ 收到消息: %d", id, len(value))
		}
	}()

	for range pingTicker.C {
		ping := Subscribe{
			RequestId: time.Now().Format("150405"),
			Action:    "ping",
			Timestamp: time.Now().UnixMilli(),
		}
		pingBytes, _ := json.Marshal(ping)
		if err := conn.WriteMessage(websocket.TextMessage, pingBytes); err != nil {
			log.Printf("[Client %d] ❌ ping 失败: %v", id, err)
			return
		}
	}
}

func main() {
	var wsURL string
	var clientCount int
	flag.StringVar(&wsURL, "url", "ws://127.0.0.1:8081/ws", "WebSocket 服务器地址")
	flag.IntVar(&clientCount, "n", 500, "客户端数量")
	flag.Parse()

	log.Printf("🚀 启动 %d 个 WebSocket 客户端...", clientCount)

	var wg sync.WaitGroup
	wg.Add(clientCount)

	for i := 0; i < clientCount; i++ {
		go newClient(i, wsURL, &wg)
		time.Sleep(time.Millisecond * 10)
	}

	wg.Wait()
}
