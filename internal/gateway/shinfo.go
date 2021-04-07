package gateway

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/model"
)

type Service interface {
	Send(topic string, userID string, message []byte) error
	Close() error
}

type Client interface {
	Send(message []byte) error
	Close() error
}

type Hub interface {
	Subscribe(topic string, client Client) error
	Unsubscribe(topic string, client Client) error
	Close() error
}

type MessageType int

const (
	basicMessage int = iota
	consumedMessage
)

type Message struct {
	RoomID      gocql.UUID
	MessageID   gocql.UUID
	UserID      gocql.UUID
	MessageType MessageType
	Entry       string
	MetaData    map[string]string
}

type UserRoom struct {
	UserID            gocql.UUID
	RoomID            gocql.UUID
	LastSeenMessageID gocql.UUID
}

type RoomService interface {
	Get(ctx context.Context, roomID gocql.UUID) (model.Room, error)
	Create(ctx context.Context, room model.Room) (model.Room, error)
	Archive(ctx context.Context, roomID gocql.UUID) error
	GetByUserID(ctx context.Context, userID gocql.UUID) ([]model.Room, error)
}

type MessageService interface {
	GetMessages(ctx context.Context, roomID gocql.UUID, limit int, lastSeenID *gocql.UUID) ([]Message, error)
	WriteMessage(ctx context.Context, roomID gocql.UUID, message Message) error
	DeleteMessage(ctx context.Context, messageID gocql.UUID) error
}

type LastSeenService interface {
	ConsumeAllMessages(ctx context.Context, userID gocql.UUID, roomID gocql.UUID) error
	GetLastSeen(ctx context.Context, userID gocql.UUID) error
}
