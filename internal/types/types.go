package types

import "time"

// Request 请求结构体.
type Request struct {
	RequestID string `json:"request_id"` // 请求唯一标识
	Timestamp int64  `json:"timestamp"`  // 发起请求的时间戳（毫秒）
	Action    string `json:"action"`     // 请求类型：subscribe | unsubscribe | ping
	Params    Params `json:"params"`     // 请求参数
}

// Params 请求参数结构体.
type Params struct {
	Topic string `json:"topic"` // 订阅主题
	Type  string `json:"type"`  // 订阅类型
}

// Response 响应结构体.
type Response struct {
	RequestID string `json:"request_id"`    // 与请求对应的 ID
	Action    string `json:"action"`        // 响应类型
	Status    string `json:"status"`        // 处理结果：ok | error
	ErrorMsg  string `json:"error_message"` // 错误信息（仅在失败时存在）
	Timestamp int64  `json:"timestamp"`     // 响应时间（毫秒）
}

// Data 推送数据结构体.
type Data struct {
	Topic   string `json:"topic"` // 推送主题
	Type    string `json:"type"`
	Payload string `json:"payload"` // 推送内容（原始 JSON 数据）
	Meta    Meta   `json:"meta"`    // 元信息
}

// Meta 元信息结构体.
type Meta struct {
	ReceiveTime time.Time `json:"receive_time"` // 消息接收时间（毫秒）
	MessageID   string    `json:"message_id"`   // 消息唯一标识
}
