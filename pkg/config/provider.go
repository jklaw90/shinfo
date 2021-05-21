package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Provider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
}

func New(fileName, path string) (Provider, error) {
	v := viper.New()
	v.SetConfigName(fileName)
	v.AddConfigPath(path)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v, nil
}
