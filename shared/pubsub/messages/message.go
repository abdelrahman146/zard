package messages

type Message interface {
	Stream() string
	Subject() string
	Consumer(group string) string
}

var Messages = []Message{&NewActivity{}}
