package pubsub

import "github.com/abdelrahman146/zard/shared/message"

type PubSub interface {
	Publish(message message.Message) error
	Subscribe(message message.Message, handler func(received []byte) error) (Subscription, error)
}

type Subscription interface {
	Unsubscribe() error
}
