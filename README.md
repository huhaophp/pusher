# pusher


pusher 是一个基于 nbio 构建的高性能 WebSocket 推送服务，支持按 topic + type 粒度的订阅管理，具备连接管理、非阻塞消息推送、后台数据源监听等能力，适用于金融、行情、消息等实时推送场景。


## 接口说明

### 请求结构（客户端 -> 服务端）

```json
{
  "request_id": "unique-id",
  "timestamp": 1720343200000,
  "action": "subscribe",
  "params": {
    "topic": "order",
    "type": "event"
  }
}

```

- action 支持：subscribe | unsubscribe | ping

### 响应结构（服务端 -> 客户端）
    
```json
{
  "request_id": "unique-id",
  "action": "subscribe",
  "status": "ok",
  "error_message": "",
  "timestamp": 1754485886585
}
```

### 推送结构（服务端 -> 客户端）

```json
{
  "topic": "order",
  "type": "event",
  "payload": "{\"content\":\"message\"}",
  "meta": {
    "receive_time": "2025-08-06T14:20:00Z",
    "message_id": "msg-001"
  }
}

```
### javascript 演示 demo
<img width="3326" height="1028" alt="image" src="https://github.com/user-attachments/assets/7ddb8d9d-b561-4284-a75c-e2ca4ec3c83c" />


