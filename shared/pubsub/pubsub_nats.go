package pubsub

import (
	"encoding/json"
	"github.com/abdelrahman146/zard/shared/logger"
	"github.com/abdelrahman146/zard/shared/message"
	"github.com/nats-io/nats.go"
	"time"
)

type natsPubSub struct {
	nc     *nats.Conn
	js     nats.JetStreamContext
	config NatsPubSubConfig
}

type NatsPubSubConfig struct {
	ResendAfter time.Duration // when a message fails to be processed, it will be resent after this duration
	Group       string
}

func NewNatsPubSub(nc *nats.Conn, config NatsPubSubConfig) PubSub {
	js, err := nc.JetStream()
	setupStreams(js)
	if err != nil {
		logger.GetLogger().Panic("failed to create jetstream context", logger.Field("error", err))
	}
	return &natsPubSub{
		nc:     nc,
		js:     js,
		config: config,
	}
}

func setupStreams(js nats.JetStreamContext) {
	messages := message.Messages
	streams := make(map[string][]string)
	for _, msg := range messages {
		if msg.Stream() == "" {
			continue
		}
		if _, ok := streams[msg.Stream()]; ok {
			streams[msg.Stream()] = append(streams[msg.Stream()], msg.Subject())
		} else {
			streams[msg.Stream()] = []string{msg.Subject()}
		}
	}
	for stream, subjects := range streams {
		_, err := js.AddStream(&nats.StreamConfig{
			Name:      stream,
			Subjects:  subjects,
			Retention: nats.WorkQueuePolicy,
			Storage:   nats.FileStorage,
		})
		if err != nil {
			logger.GetLogger().Panic("failed to create stream", logger.Field("error", err))
		}
	}
}

func (n *natsPubSub) Publish(message message.Message) error {
	data, _ := json.Marshal(message)
	switch {
	case message.Stream() == "":
		return n.nc.Publish(message.Subject(), data)
	default:
		_, err := n.js.Publish(message.Subject(), data)
		return err
	}
}

func (n *natsPubSub) Subscribe(message message.Message, handler func(received []byte) error) (Subscription, error) {
	consumer := message.Consumer(n.config.Group)
	switch {
	case message.Stream() == "":
		sub, err := n.nc.QueueSubscribe(message.Subject(), consumer, func(msg *nats.Msg) {
			if err := handler(msg.Data); err != nil {
				_ = msg.Nak()
			} else {
				_ = msg.Ack()
			}
		})
		return sub, err
	default:
		sub, err := n.js.QueueSubscribe(message.Subject(), consumer, func(natsMsg *nats.Msg) {
			if err := handler(natsMsg.Data); err != nil {
				_ = natsMsg.NakWithDelay(n.config.ResendAfter)
			} else {
				_ = natsMsg.Ack()
			}
		}, nats.ManualAck(), nats.Durable(consumer), nats.DeliverAll())
		return sub, err
	}
}
