package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Session struct {
	UserID   uuid.UUID
	DeviceID uuid.UUID
	ServerID uuid.UUID
}

func (s *Session) Key() string {
	return fmt.Sprintf("%s:%s", s.UserID.String(), s.DeviceID.String())
}
