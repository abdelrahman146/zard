package provider

import (
	"github.com/abdelrahman146/zard/shared/logger"
	capi "github.com/hashicorp/consul/api"
)

type ConsulProvider interface {
	GetClient() *capi.Client
	Close()
}

type consulProvider struct {
	client *capi.Client
}

func InitConsulProvider(address string) ConsulProvider {
	client, err := capi.NewClient(&capi.Config{Address: address})
	if err != nil {
		logger.GetLogger().Panic("Failed to connect to consulProvider", logger.Field("error", err))
	}
	return &consulProvider{
		client: client,
	}
}

func (c *consulProvider) GetClient() *capi.Client {
	return c.client
}

func (c *consulProvider) Close() {
	c.client = nil
}
