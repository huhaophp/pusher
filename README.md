# pusher

pusher 是一个基于 nbio 构建的高性能 WebSocket 推送服务，支持按 topic + type 粒度的订阅管理，具备连接管理、非阻塞消息推送、后台数据源监听等能力，适用于金融、行情、消息等实时推送场景。

## 功能特性

- 高性能 WebSocket 服务，基于 nbio 实现
- 支持按 topic + type 维度的细粒度订阅管理
- 多数据源支持：Redis 和 Kafka
- 非阻塞消息推送机制
- 连接监控和管理

## 项目目录

```text
pusher/
├── client/                 # 测试客户端
├── config/                 # 配置文件和加载逻辑
├── internal/
│   ├── source/            # 数据源实现（Redis/Kafka）
│   ├── types/             # 数据结构定义
│   └── ws/                # WebSocket 相关实现
├── pkg/
│   ├── logger/            # 日志模块
│   ├── redis/             # Redis 客户端初始化
│   ├── kafka/             # Kafka 客户端初始化
│   └── utils/             # 工具函数
├── main.go                # 程序入口
└── go.mod
```

## 核心组件

### 订阅管理器 (SubscriptionManager)

- 管理所有主题订阅
- 支持 Redis 和 Kafka 两种数据源
- 提供订阅和取消订阅接口

### 主题订阅 (TopicSubscription)

- 每个主题对应一个实例
- 管理该主题下所有类型的订阅者
- 从数据源拉取消息并推送给订阅者

### 订阅者 (Subscriber)

- 封装 WebSocket 连接
- 使用 channel 实现非阻塞消息推送
- 自动处理连接异常和关闭

### 数据源 (Source)

- 定义统一的数据源接口
- Redis 实现：基于 Redis Pub/Sub
- Kafka 实现：基于 Kafka Consumer

### 启动方式

```bash
# 编译
go build -o pusher main.go

# 运行
./pusher

# 或直接运行
go run main.go

```

### 数据流说明

1. 客户端通过 WebSocket 连接到服务
2. 客户端发送 subscribe 消息订阅指定 topic 和 type 
3. 服务端通过 Redis 或 Kafka 监听对应 topic 的消息 
4. 当有新消息时，服务端根据订阅信息推送给对应客户端 
5. 客户端可以通过 unsubscribe 取消订阅 
6. 客户端定期发送 ping 消息保持连接

## 接口说明

### subscribe （客户端 -> 服务端）

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

### subscribe 响应结构（服务端 -> 客户端）

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

### unsubscribe （客户端 -> 服务端）

```json
{
  "request_id": "unique-id",
  "action": "unsubscribe",
  "timestamp": 123123123,
  "params": {
    "topic": "topic1",
    "type": "test"
  }
}
```

### unsubscribe 响应结构（服务端 -> 客户端）

```json
{
  "request_id": "1234",
  "action": "unsubscribe",
  "status": "ok",
  "error_message": "",
  "timestamp": 1754926705364
}
```

### ping （客户端 -> 服务端）

```json
{
  "request_id": "1234",
  "action": "ping",
  "timestamp": 123123123
}
```

### ping 响应结构（服务端 -> 客户端）

```json
{
  "request_id": "1234",
  "action": "ping",
  "status": "ok",
  "error_message": "",
  "timestamp": 1754926792165
}
```

### 演示 demo.html

<img width="3326" height="1028" alt="image" src="https://github.com/user-attachments/assets/7ddb8d9d-b561-4284-a75c-e2ca4ec3c83c" />
