package database

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/model"
)

const (
	userExpire = 10 * time.Minute
)

type RedisCachce struct {
	conn *redis.Client
}

func NewCache(conn *redis.Client) *RedisCachce {
	return &RedisCachce{
		conn: conn,
	}
}

func (c *RedisCachce) SetRooms(userID gocql.UUID, roomsLastUpdated []model.RoomLastUpdated) error {
	rooms := []redis.Z{}
	for _, rlu := range roomsLastUpdated {
		rooms = append(rooms, redis.Z{
			Score:  float64(rlu.LastUpdated.Time().Unix()),
			Member: rlu.RoomID.String(),
		})
	}
	cmd := c.conn.ZAdd(key(userID), rooms...)
	c.conn.Expire(key(userID), userExpire)
	return cmd.Err()
}

func (c *RedisCachce) GetRooms(userID gocql.UUID, limit, offset int64) ([]gocql.UUID, error) {
	tmpRoomsStr := []string{}
	cmd := c.conn.ZRange(key(userID), offset, offset+limit)
	if err := cmd.ScanSlice(&tmpRoomsStr); err != nil {
		return nil, err
	}

	rooms := []gocql.UUID{}
	for _, rs := range tmpRoomsStr {
		r, err := gocql.ParseUUID(rs)
		if err != nil {
			continue
		}
		rooms = append(rooms, r)
	}
	c.conn.Expire(key(userID), userExpire)

	return rooms, nil
}

func (c *RedisCachce) UpdateRoom(userID gocql.UUID, room model.RoomLastUpdated) error {
	cmd := c.conn.ZAddXX(key(userID), redis.Z{
		Score:  float64(room.LastUpdated.Time().Unix()),
		Member: room.RoomID.String(),
	})
	return cmd.Err()
}

func key(userID gocql.UUID) string {
	return fmt.Sprintf("users:%s:rooms", userID.String())
}
