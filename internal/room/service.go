package room

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	pb "github.com/jklaw90/shinfo/internal/pb/room"
	"github.com/jklaw90/shinfo/internal/room/database"
	"github.com/jklaw90/shinfo/pkg/config"
	"github.com/jklaw90/shinfo/pkg/model"
	"google.golang.org/grpc"
)

type RoomService struct {
	repo  *database.RoomRepo
	cache *database.RedisCachce
}

var _ Service = (*RoomService)(nil)

func NewService(session *gocql.Session, rClient *redis.Client) *RoomService {
	return &RoomService{
		repo:  database.NewDB(session),
		cache: database.NewCache(rClient),
	}
}

func NewServer(
	ctx context.Context,
	cfg config.Provider,
	session *gocql.Session,
	rClient *redis.Client,
) (func() error, error) {

	service := &RoomService{
		repo:  database.NewDB(session),
		cache: database.NewCache(rClient),
	}

	lis, err := net.Listen("tcp", cfg.GetString("room.address"))
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	server, err := NewRoomServer(service)
	if err != nil {
		return nil, err
	}

	pb.RegisterRoomServer(s, server)

	return func() error {
		return s.Serve(lis)
	}, nil
}

func (s *RoomService) Get(ctx context.Context, roomID gocql.UUID) (model.Room, error) {
	r, err := s.repo.Get(ctx, roomID)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (s *RoomService) Create(ctx context.Context, params model.RoomCreate) (model.Room, error) {
	now := time.Now().UTC()
	room := model.Room{
		ID:        gocql.UUIDFromTime(now),
		Name:      params.Name,
		Type:      params.Type,
		Public:    params.Public,
		Archived:  params.Public,
		CreatedAt: now,
	}
	if err := s.repo.Create(ctx, room); err != nil {
		return room, err
	}
	return room, nil
}

func (s *RoomService) AddUser(ctx context.Context, roomID gocql.UUID, user model.User) error {
	room, err := s.Get(ctx, roomID)
	if err != nil {
		return err
	}

	for _, u := range room.Users {
		if u.ID == user.ID {
			return errors.New("user already exists")
		}
	}

	if err := s.repo.AddUser(ctx, room, user); err != nil {
		return err
	}

	return nil
}

func (s *RoomService) GetByUserID(ctx context.Context, userID gocql.UUID, limit int64, lastSeen *gocql.UUID) (model.UserRooms, error) {
	var resp model.UserRooms
	if limit == 0 || limit > 100 {
		limit = 25
	}

	rooms, err := s.repo.GetByUserID(ctx, userID, limit, lastSeen)
	if err != nil {
		return resp, err
	}
	resp.Rooms = rooms
	if len(rooms) == int(limit) {
		resp.NextID = &rooms[len(rooms)-1].RoomID
	}

	return resp, nil
}

func (s *RoomService) GetUsers(ctx context.Context, roomID gocql.UUID, limit int64, lastSeen *gocql.UUID) (model.RoomUsers, error) {
	var roomUsers model.RoomUsers

	users, err := s.repo.GetUsersByRoomID(ctx, roomID, limit, lastSeen)
	if err != nil {
		return roomUsers, err
	}

	for _, ru := range users {
		roomUsers.Users = append(roomUsers.Users, model.RoomUser{
			User:     ru.User,
			JoinedAt: ru.JoinedAt.Time(),
		})
	}
	if len(users) == int(limit) {
		roomUsers.NextID = &users[limit-1].JoinedAt
	}

	return roomUsers, nil
}

func (s *RoomService) OnMessage(ctx context.Context, msg model.Message) error {
	r, err := s.repo.Get(ctx, msg.RoomID)
	if err != nil {
		return err
	}

	for _, uid := range r.Users {
		s.cache.UpdateRoom(uid.ID, model.RoomLastUpdated{
			RoomID:      msg.RoomID,
			LastUpdated: msg.ID,
		})
	}

	return nil
}
