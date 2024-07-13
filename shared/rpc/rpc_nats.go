package rpc

import (
	"encoding/json"
	"github.com/abdelrahman146/zard/shared/provider"
	"github.com/abdelrahman146/zard/shared/rpc/requests"
	"github.com/abdelrahman146/zard/shared/validator"
	"github.com/nats-io/nats.go"
	"time"
)

type natsRPC struct {
	nc     *nats.Conn
	v      validator.Validator
	config NatsRPCConfig
}

type NatsRPCConfig struct {
	Timeout time.Duration
	Group   string
}

func NewNatsRPC(nts *provider.natsProvider, v validator.Validator, config NatsRPCConfig) RPC {
	return &natsRPC{
		nc:     nts.GetConn(),
		v:      v,
		config: config,
	}
}

func (n *natsRPC) Request(req requests.Request) (resp []byte, err error) {
	if err = n.v.ValidateStruct(req); err != nil {
		return nil, err
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	msg, err := n.nc.Request(req.Subject(), data, n.config.Timeout)
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}

func (n *natsRPC) Handle(req requests.Request, handler func(req []byte) (resp []byte)) error {
	_, err := n.nc.QueueSubscribe(req.Subject(), req.Consumer(n.config.Group), func(msg *nats.Msg) {
		resp := handler(msg.Data)
		_ = msg.Respond(resp)
	})

	return err
}
