package gateway

import (
	"github.com/go-redis/redis"
)

type RedisHub struct {
	rdb    *redis.Client
	pubsub *redis.PubSub
	done   chan bool
	BasicHub
}

func NewRedisHub(redisClient *redis.Client) *RedisHub {
	h := &RedisHub{
		rdb:    redisClient,
		pubsub: redisClient.Subscribe(),
		done:   make(chan bool),
	}
	return h
}

func (h *RedisHub) Publish(topic string, message []byte) error {
	h.rdb.Publish(topic, message)
	return nil
}

func (h *RedisHub) Listen() {
	stream := h.pubsub.Channel()
	for {
		select {
		case message := <-stream:
			h.Broadcast(message.Channel, []byte(message.Payload))
		case <-h.done:
			h.pubsub.Close()
			return
		}
	}
}

func (h *RedisHub) Subscribe(topic string, client Client) error {
	h.BaseSubscribe(topic, client, func(topic string) error {
		return h.pubsub.Subscribe(topic)
	})
	return nil
}

func (h *RedisHub) Unsubscribe(topic string, client Client) error {
	h.BaseUnsubscribe(topic, client, func(topic string) error {
		return h.pubsub.Unsubscribe(topic)
	})
	return nil
}

func (h *RedisHub) Close() error {
	h.BaseClose()
	h.done <- true
	return nil
}
