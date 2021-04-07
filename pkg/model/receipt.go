package model

import (
	"time"

	"github.com/gocql/gocql"
)

type ReadReceipt struct {
	RoomID      gocql.UUID `cql:"room_id"`
	UserID      gocql.UUID `cql:"user_id"`
	UnreadCount int        `cql:"unread_count"`
	UpdatedAt   time.Time  `cql:"updated_at"`
}
