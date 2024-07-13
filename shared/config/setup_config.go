package config

import (
	"encoding/json"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/provider"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func GetEnvConfig(conf Config) {
	_ = godotenv.Load()
	envs := make(map[string]interface{})
	prefix := "ZARD_"
	for _, key := range os.Environ() {
		if strings.HasPrefix(key, prefix) {
			pair := strings.SplitN(key, "=", 2)
			pair[0] = strings.TrimPrefix(pair[0], prefix)
			envs[pair[0]] = shared.Utils.Strings.Parse(pair[1])
		}
	}
	setBulk(conf, "env", envs)
}

func GetConsulConfig(c *provider.consulProvider, kvPath string, conf Config) error {
	client := c.GetClient()
	kv := client.KV()
	pair, _, err := kv.Get(kvPath, nil)
	if err != nil {
		return err
	}
	var configData map[string]interface{}
	err = json.Unmarshal(pair.Value, &configData)
	if err != nil {
		return err
	}
	setBulk(conf, "consul", configData)
	return nil
}
