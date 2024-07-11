package config

import (
	"github.com/spf13/viper"
)

type viperConfig struct {
	v *viper.Viper
}

func NewViperConfig() Config {
	v := viper.New()
	v.AutomaticEnv()
	return &viperConfig{v: v}
}

func (v *viperConfig) GetString(key string) string {
	return v.v.GetString(key)
}

func (v *viperConfig) GetInt(key string) int {
	return v.v.GetInt(key)
}

func (v *viperConfig) GetFloat(key string) float64 {
	return v.v.GetFloat64(key)
}

func (v *viperConfig) GetBool(key string) bool {
	return v.v.GetBool(key)
}

func (v *viperConfig) Set(key string, value interface{}) {
	v.v.Set(key, value)
}
