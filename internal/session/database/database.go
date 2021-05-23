package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jklaw90/shinfo/internal/session/model"
)

const (
	SessionTTL = time.Hour * 1
)

type Db struct {
	client *redis.Client
}

func New(ctx context.Context) (*Db, error) {
	return nil, nil
}

func (d *Db) SessionSet(ctx context.Context, session model.Session) error {
	return d.client.Set(session.Key(), session.ServerID, SessionTTL).Err()
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
