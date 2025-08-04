package model

const (
	StatusOK   = "ok"   // 成功状态
	StatusFail = "fail" // 失败状态
)

const (
	ActionSubscribe   = "subscribe"   // 订阅
	ActionUnsubscribe = "unsubscribe" // 取消订阅
	ActionPing        = "ping"        // 心跳
)

var (
	AllAction = []string{ActionSubscribe, ActionUnsubscribe, ActionPing}
)
