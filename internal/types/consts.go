package types

import "time"

var (
	AllAction = []string{ActionSubscribe, ActionUnsubscribe, ActionPing}
)

const (
	StatusOK   = "ok"   // 成功状态
	StatusFail = "fail" // 失败状态

	ActionSubscribe   = "subscribe"   // 订阅
	ActionUnsubscribe = "unsubscribe" // 取消订阅
	ActionPing        = "ping"        // 心跳

	SubscriberSendChanSize = 100 // 订阅者发送chan大小
	ConnDeadlineTime       = time.Minute * 2
)
