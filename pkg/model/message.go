package model

import (
	"github.com/gocql/gocql"
)

type MessageType int64

const (
	ChatMessage MessageType = iota
	SeenMessage
	TypingMessage
)

type Message struct {
	ID       gocql.UUID        `cql:"id"`
	RoomID   gocql.UUID        `cql:"room_id"`
	User     User              `cql:"user,omitempty"`
	Type     *MessageType      `cql:"message_type"`
	Entry    string            `cql:"entry,omitempty"`
	Metadata map[string]string `cql:"metadata,omitempty"`
}

type MessageCreate struct {
	RoomID   gocql.UUID        `cql:"room_id"`
	User     User              `cql:"user,omitempty"`
	Type     *MessageType      `cql:"message_type"`
	Entry    string            `cql:"entry,omitempty"`
	Metadata map[string]string `cql:"metadata,omitempty"`
}

type MessageList struct {
	Messages []Message
	NextID   *gocql.UUID
}
