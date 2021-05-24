package config

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("config", "../../config/local.yaml", "config file filepath")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

type Provider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	GetDuration(key string) time.Duration
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
}

func New() (Provider, error) {
	v := viper.New()
	v.SetConfigFile(viper.GetString("config"))
	v.WatchConfig()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}
