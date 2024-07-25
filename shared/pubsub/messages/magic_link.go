package messages

import "time"

type MagicLinkMessage struct {
	UserID    string    `json:"userId"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	WithReset bool      `json:"withReset"`
	Timestamp time.Time `json:"timestamp"`
}

func (m *MagicLinkMessage) Stream() string {
	return "account"
}

func (m *MagicLinkMessage) Subject() string {
	return "account.magiclink.send"
}

func (m *MagicLinkMessage) Consumer(group string) string {
	if group != "" {
		return "magiclink_send_" + group
	}
	return "magiclink_send"
}
