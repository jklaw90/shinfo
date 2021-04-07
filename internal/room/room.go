package room

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/model"
	usermodel "github.com/jklaw90/shinfo/pkg/model"
)

type Service interface {
	Get(ctx context.Context, roomID gocql.UUID) (model.Room, error)
	GetByUserID(ctx context.Context, userID gocql.UUID, limit int64, lastSeen *gocql.UUID) (model.UserRooms, error)
	GetUsers(ctx context.Context, roomID gocql.UUID, limit int64, lastSeen *gocql.UUID) (model.RoomUsers, error)
	Create(ctx context.Context, room model.RoomCreate) (model.Room, error)
	AddUser(ctx context.Context, roomID gocql.UUID, user usermodel.User) error
	OnMessage(ctx context.Context, msg model.Message) error
}
