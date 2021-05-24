package database

import (
	"github.com/go-redis/redis"
	"github.com/jklaw90/shinfo/pkg/config"
)

func NewRedis(cfg config.Provider) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetString("redis.address"),
		DB:       cfg.GetInt("redis.database"),
		Password: cfg.GetString("redis.password"),
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return client, nil
}
