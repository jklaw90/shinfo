package message

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/internal/message/database"
	"github.com/jklaw90/shinfo/pkg/model"
)

type MessageService struct {
	repo  *database.MessageRepo
	cache *database.RedisCachce
}

func NewService(session *gocql.Session, redisClient *redis.Client) *MessageService {
	return &MessageService{
		repo:  database.NewDB(session),
		cache: database.NewCache(redisClient),
	}
}

type ListParams struct {
	RoomID gocql.UUID
	Limit  int64
	NextID *gocql.UUID
}

func (s *MessageService) Get(ctx context.Context, roomID, messageID gocql.UUID) (model.Message, error) {
	msg, err := s.repo.Get(ctx, roomID, messageID)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

func (s *MessageService) List(ctx context.Context, params ListParams) (model.MessageList, error) {
	var msgList model.MessageList
	msgs, err := s.getMessages(ctx, params)
	if err != nil {
		return msgList, err
	}
	msgList.Messages = msgs
	if len(msgs) == int(params.Limit) {
		msgList.NextID = &msgs[len(msgs)-1].ID
	}
	return msgList, nil
}

func (s *MessageService) getMessages(ctx context.Context, params ListParams) ([]model.Message, error) {
	var msgs []model.Message
	if params.NextID == nil {
		msgs, err := s.cache.GetRoomMessages(params.RoomID, params.Limit)
		if err != nil {
			return msgs, err
		}
		if len(msgs) > 0 {
			return msgs, nil
		}
	}
	msgs, err := s.repo.GetMessages(ctx, params.RoomID, params.Limit, params.NextID)
	if err != nil {
		return msgs, err
	}
	s.cache.SetRoomMessages(params.RoomID, msgs...)
	return msgs, nil
}

func (s *MessageService) Delete(ctx context.Context, roomID, messageID gocql.UUID) error {
	err := s.repo.Delete(ctx, roomID, messageID)
	if err != nil {
		return err
	}
	return nil
}

func (s *MessageService) AddMessage(ctx context.Context, m model.MessageCreate) (model.Message, error) {
	msg := model.Message{
		ID:       gocql.TimeUUID(),
		RoomID:   m.RoomID,
		User:     m.User,
		Type:     m.Type,
		Entry:    m.Entry,
		Metadata: m.Metadata,
	}
	if err := s.repo.Create(ctx, msg); err != nil {
		return msg, err
	}
	s.cache.SetRoomMessages(m.RoomID, msg)
	return msg, nil
}
