package messages

import "time"

type AuthOTPCreated struct {
	Value     string        `json:"value"`
	Target    string        `json:"target"`
	Otp       string        `json:"otp"`
	Ttl       time.Duration `json:"ttl"`
	Timestamp time.Time     `json:"timestamp"`
}

func (a *AuthOTPCreated) Stream() string {
	return "account"
}

func (a *AuthOTPCreated) Subject() string {
	return "account.auth.otp.created"
}

func (a *AuthOTPCreated) Consumer(group string) string {
	if group != "" {
		return "account_auth_otp_created_" + group
	}
	return "account_auth_otp_created"
}
