package cache

import (
	"github.com/abdelrahman146/zard/shared/logger"
	"github.com/abdelrahman146/zard/shared/provider"
	"github.com/nats-io/nats.go"
	"strings"
)

type natsCache struct {
	bucket nats.KeyValue
	config *nats.KeyValueConfig
}

func NewNatsCache(nts *provider.natsProvider, config *nats.KeyValueConfig) Cache {
	c := &natsCache{
		config: config,
	}
	if err := c.init(nts.GetConn()); err != nil {
		logger.GetLogger().Panic("failed to initialize nats cache", logger.Field("error", err))
	}
	return c
}

func (c *natsCache) init(nc *nats.Conn) error {
	js, err := nc.JetStream()
	if err != nil {
		return err
	}
	bucket, err := js.KeyValue(c.config.Bucket)
	if err != nil {
		bucket, err = js.CreateKeyValue(c.config)
		if err != nil {
			return err
		}
	}
	c.bucket = bucket
	return nil
}

func (c *natsCache) Get(keyPath []string) (value []byte, err error) {
	key := strings.Join(keyPath, ".")
	entry, err := c.bucket.Get(key)
	if err != nil {
		return nil, err
	}
	return entry.Value(), nil
}

func (c *natsCache) Set(keyPath []string, value []byte) error {
	key := strings.Join(keyPath, ".")
	_, err := c.bucket.Put(key, value)
	return err
}

func (c *natsCache) Delete(keyPath []string) error {
	key := strings.Join(keyPath, ".")
	return c.bucket.Delete(key)
}

func (c *natsCache) Keys() (keys []string, err error) {
	return c.bucket.Keys()
}
