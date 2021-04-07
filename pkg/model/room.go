package model

import (
	"time"

	"github.com/gocql/gocql"
)

type RoomType int64

const (
	// Classic room has users and doesn't
	Classic RoomType = iota
	// Ephemeral room doesn't save messages
	Ephemeral
)

type Room struct {
	ID        gocql.UUID `cql:"id"`
	Name      string     `cql:"name,omitempty"`
	Users     []User     `cql:"users"`
	Type      *RoomType  `cql:"type"`
	Public    *bool      `cql:"public"`
	Archived  *bool      `cql:"archived"`
	CreatedAt time.Time  `cql:"created_at"`
}

type RoomUser struct {
	User     User
	JoinedAt time.Time
}

type RoomUsers struct {
	Users  []RoomUser
	NextID *gocql.UUID
}

type UserByRoom struct {
	RoomID   gocql.UUID `cql:"room_id"`
	JoinedAt gocql.UUID `cql:"joined_at"`
	User     User
}

type RoomByUser struct {
	RoomID    gocql.UUID `cql:"room_id"`
	RoomName  string     `cql:"room_name,omitempty"`
	Type      *RoomType  `cql:"type"`
	Public    *bool      `cql:"public"`
	Archived  *bool      `cql:"archived"`
	CreatedAt time.Time  `cql:"created_at"`
	JoinedAt  time.Time  `cql:"joined_at"`
}

type UserRooms struct {
	Rooms  []RoomByUser
	NextID *gocql.UUID
}

type RoomLastUpdated struct {
	RoomID      gocql.UUID `cql:"room_id"`
	LastUpdated gocql.UUID `cql:"updated_at"`
}

type RoomCreate struct {
	Name     string    `cql:"name,omitempty"`
	Type     *RoomType `cql:"type"`
	Public   *bool     `cql:"public"`
	Archived *bool     `cql:"archived"`
}
