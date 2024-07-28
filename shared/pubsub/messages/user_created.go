package messages

import "time"

type UserCreatedMessage struct {
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone"`
	Timestamp time.Time `json:"timestamp"`
}

func (m *UserCreatedMessage) Stream() string {
	return "account"
}

func (m *UserCreatedMessage) Subject() string {
	return "account.user.created"
}

func (m *UserCreatedMessage) Consumer(group string) string {
	if group != "" {
		return "account_user_created_" + group
	}
	return "account_user_created"
}
