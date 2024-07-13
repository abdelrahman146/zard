package config

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetFloat(key string) float64
	GetBool(key string) bool
	Set(key string, value interface{})
}

func setBulk(conf Config, prefix string, data map[string]interface{}) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		switch v := value.(type) {
		case map[string]interface{}:
			setBulk(conf, fullKey, v)
		default:
			conf.Set(fullKey, v)
		}
	}
}
