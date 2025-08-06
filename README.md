# pusher


pusher 是一个基于 nbio 构建的高性能 WebSocket 推送服务，支持按 topic + type 粒度的订阅管理，具备连接管理、非阻塞消息推送、后台数据源监听等能力，适用于金融、行情、消息等实时推送场景。


## 接口说明

### 请求结构（客户端 -> 服务端）

```json
{
  "request_id": "abc123",
  "timestamp": 1720343200000,
  "action": "subscribe",
  "params": {
    "topic": "spot.kline.btcusdt",
    "type": "1m"
  }
}

```

- action 支持：subscribe | unsubscribe | ping

### 响应结构（服务端 -> 客户端）
    
```json
{
  "request_id": "abc123",
  "action": "subscribe",
  "status": "ok",
  "error_message": "",
  "timestamp": 1754485886585
}
```

### 推送结构（服务端 -> 客户端）

```json
{
  "topic": "spot.kline.btcusdt",
  "type": "1m",
  "payload": "{\"open\":123.45,\"high\":125.6,...}",
  "meta": {
    "receive_time": "2025-08-06T14:20:00Z",
    "message_id": "msg-001"
  }
}

```

