package message

import (
	"time"
)

type NewActivity struct {
	AppID          string    `json:"appId"`
	SubscriptionID string    `json:"subscriptionId"`
	Timestamp      time.Time `json:"timestamp"`
	Service        string    `json:"service"`
	Action         string    `json:"action"`
	Source         string    `json:"source"`
	RequestID      string    `json:"requestId"`
}

func (m *NewActivity) Stream() string {
	return "subscription"
}

func (m *NewActivity) Subject() string {
	return "sub.activity.new"
}

func (m *NewActivity) Consumer(group string) string {
	if group != "" {
		return "activity_new_" + group
	}
	return "activity_new"
}
