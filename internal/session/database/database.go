package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jklaw90/shinfo/internal/session/model"
	"github.com/jklaw90/shinfo/pkg/config"
)

const (
	SessionDefaultTTL = time.Minute * 10
)

type Db struct {
	client *redis.Client
	cfg    config.Provider
}

func New(ctx context.Context, cfg config.Provider) (*Db, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetString("redis.address"),
		DB:       cfg.GetInt("redis.database"),
		Password: cfg.GetString("redis.password"),
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &Db{
		client: client,
		cfg:    cfg,
	}, nil
}

func (d *Db) SessionSet(ctx context.Context, session model.Session) error {
	return d.client.Set(session.Key(), session.ServerID, d.getTTL()).Err()
}

func (d *Db) SessionDelete(ctx context.Context, session model.Session) error {
	return d.client.Del(session.Key()).Err()
}

func (d *Db) SessionsGetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Session, error) {
	var sessions []model.Session

	cmd := d.client.Get(fmt.Sprintf("%s:*", userID.String()))
	if err := cmd.Scan(&sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (d *Db) getTTL() time.Duration {
	if ttl := d.cfg.GetDuration("session.ttl"); ttl != 0 {
		return ttl
	}
	return SessionDefaultTTL
}
