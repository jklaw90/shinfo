package session

import (
	"context"
	"sync"

	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/model"
)

type Session struct {
	User   model.User
	RoomID gocql.UUID
	Client string
}

type SessionService struct {
	sessions map[string]Session
	rooms    map[gocql.UUID][]string
	sync.RWMutex
}

func (s *SessionService) MarkActive(ctx context.Context, session Session) error {
	s.Lock()
	defer s.Unlock()

	return nil
}

func (s *SessionService) MarkInActive(ctx context.Context, session Session) error {
	s.Lock()
	defer s.Unlock()

	return nil
}

func (s *SessionService) OnMessage(ctx context.Context, msg model.Message) error {
	s.RLock()
	defer s.RUnlock()

	return nil
}
