package pubsub

import (
	"github.com/abdelrahman146/zard/shared/pubsub/messages"
)

type PubSub interface {
	Publish(message messages.Message) error
	Subscribe(message messages.Message, handler func(received []byte) error) (Subscription, error)
}

type Subscription interface {
	Unsubscribe() error
}
