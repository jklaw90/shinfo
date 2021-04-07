package database

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/model"
	"github.com/jklaw90/shinfo/pkg/utils"
)

const (
	chatExpire = 10 * time.Minute

	defaultLimit = 50
)

type RedisCachce struct {
	conn *redis.Client
}

func NewCache(conn *redis.Client) *RedisCachce {
	return &RedisCachce{
		conn: conn,
	}
}

func (c *RedisCachce) GetRoomMessages(roomID gocql.UUID, limit int64) ([]model.Message, error) {
	msgs := []model.Message{}
	key := key(roomID)

	resp := c.conn.LRange(key, 0, limit)
	if err := resp.ScanSlice(&msgs); err != nil {
		return msgs, err
	}

	return msgs, nil
}

func (c *RedisCachce) SetRoomMessages(roomID gocql.UUID, msgs ...model.Message) {
	key := key(roomID)
	c.conn.LPush(key, msgs)
	c.conn.Expire(key, chatExpire)
	if utils.RandomInRange(5)%5 == 0 {
		c.conn.LTrim(key, 0, defaultLimit)
	}
}

func key(roomID gocql.UUID) string {
	return fmt.Sprintf("rooms:%s:messages", roomID.String())
}
