package requests

type Request interface {
	Subject() string
	Consumer(group string) string
}

var Requests = []Request{&GetUserRequest{}}
