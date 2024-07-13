package provider

import (
	"github.com/abdelrahman146/zard/shared/logger"
	"github.com/nats-io/nats.go"
)

type NatsProvider interface {
	GetConn() *nats.Conn
	GetJs() nats.JetStreamContext
	Close()
}

type natsProvider struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func InitNatsProvider(address string) NatsProvider {
	nc, err := nats.Connect(address)
	if err != nil {
		logger.GetLogger().Panic("Failed to connect to NATS", logger.Field("error", err))
	}
	js, err := nc.JetStream()
	if err != nil {
		logger.GetLogger().Panic("Failed to connect to JetStream", logger.Field("error", err))
	}
	return &natsProvider{
		nc: nc,
		js: js,
	}
}

func (n *natsProvider) GetConn() *nats.Conn {
	return n.nc
}

func (n *natsProvider) GetJs() nats.JetStreamContext {
	return n.js
}

func (n *natsProvider) Close() {
	n.nc.Close()
	logger.GetLogger().Info("NATS connection closed")
}
