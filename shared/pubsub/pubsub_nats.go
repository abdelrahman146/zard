package pubsub

import (
	"encoding/json"
	"github.com/abdelrahman146/zard/shared/logger"
	"github.com/abdelrahman146/zard/shared/provider"
	"github.com/abdelrahman146/zard/shared/pubsub/messages"
	"github.com/nats-io/nats.go"
	"time"
)

type natsPubSub struct {
	nts    provider.NatsProvider
	config NatsPubSubConfig
}

type NatsPubSubConfig struct {
	ResendAfter time.Duration // when a comm fails to be processed, it will be resent after this duration
	Group       string
}

func NewNatsPubSub(nts provider.NatsProvider, config NatsPubSubConfig) PubSub {
	setupStreams(nts.GetJs())
	return &natsPubSub{
		nts:    nts,
		config: config,
	}
}

func setupStreams(js nats.JetStreamContext) {
	msgs := messages.Messages
	streams := make(map[string][]string)
	for _, msg := range msgs {
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

func (p *natsPubSub) Publish(message messages.Message) error {
	data, _ := json.Marshal(message)
	switch {
	case message.Stream() == "":
		return p.nts.GetConn().Publish(message.Subject(), data)
	default:
		_, err := p.nts.GetJs().Publish(message.Subject(), data)
		return err
	}
}

func (p *natsPubSub) Subscribe(message messages.Message, handler func(received []byte) error) (Subscription, error) {
	consumer := message.Consumer(p.config.Group)
	switch {
	case message.Stream() == "":
		sub, err := p.nts.GetConn().QueueSubscribe(message.Subject(), consumer, func(msg *nats.Msg) {
			if err := handler(msg.Data); err != nil {
				_ = msg.Nak()
			} else {
				_ = msg.Ack()
			}
		})
		return sub, err
	default:
		sub, err := p.nts.GetJs().QueueSubscribe(message.Subject(), consumer, func(natsMsg *nats.Msg) {
			if err := handler(natsMsg.Data); err != nil {
				_ = natsMsg.NakWithDelay(p.config.ResendAfter)
			} else {
				_ = natsMsg.Ack()
			}
		}, nats.ManualAck(), nats.Durable(consumer), nats.DeliverAll())
		return sub, err
	}
}
