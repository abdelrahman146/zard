package config

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetFloat(key string) float64
	GetBool(key string) bool
	Set(key string, value interface{})
}
