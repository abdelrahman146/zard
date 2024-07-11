package config

import (
	"encoding/json"
	"github.com/abdelrahman146/zard/shared/logger"
	capi "github.com/hashicorp/consul/api"
)

func NewConsulConfig(address string) Config {
	client, err := capi.NewClient(&capi.Config{Address: address})
	if err != nil {
		logger.GetLogger().Panic("failed to create consul client", logger.Field("error", err))
	}
	kv := client.KV()
	pair, _, err := kv.Get("app", nil)
	if err != nil {
		logger.GetLogger().Panic("failed to get config from consul", logger.Field("error", err))
	}
	var configData map[string]interface{}
	err = json.Unmarshal(pair.Value, &configData)
	if err != nil {
		logger.GetLogger().Panic("failed to unmarshal config data", logger.Field("error", err))
	}
	viperConfig := NewViperConfig()
	setConfig(viperConfig, "", configData)
	return viperConfig
}

func setConfig(conf Config, prefix string, data map[string]interface{}) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		switch v := value.(type) {
		case map[string]interface{}:
			setConfig(conf, fullKey, v)
		default:
			conf.Set(fullKey, v)
		}
	}
}
