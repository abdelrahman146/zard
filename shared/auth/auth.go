package auth

import (
	"errors"
)

var (
	InvalidCredentialsTypeError = errors.New("invalid credential type")
	InvalidCredentialsError     = errors.New("invalid credentials")
	InvalidTokenError           = errors.New("invalid token")
	UnauthorizedError           = errors.New("unauthorized")
	TokenExpiredError           = errors.New("token expired")
)

type Auth interface {
	Create(subject string) (token string, err error)
	Authenticate(token string) (claims Claims, err error)
	Revoke(token string) error
}
