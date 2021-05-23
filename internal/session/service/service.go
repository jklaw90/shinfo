package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jklaw90/shinfo/internal/session/database"
	"github.com/jklaw90/shinfo/internal/session/model"
	"github.com/jklaw90/shinfo/pkg/config"
	"github.com/jklaw90/shinfo/pkg/logging"
)

type SessionService struct {
	db  *database.Db
	cfg *config.Provider
}

func New(ctx context.Context, cfg *config.Provider) *SessionService {
	return &SessionService{
		cfg: cfg,
	}
}

func (s *SessionService) Register(ctx context.Context, session model.Session) error {
	logger := logging.FromContext(ctx)

	logger.Debugw("registering sessoin", "userID", session.UserID, "deviceID", session.DeviceID, "serverID", session.ServerID)

	if err := s.db.SessionSet(ctx, session); err != nil {
		return err
	}

	return nil
}

func (s *SessionService) Unregister(ctx context.Context, session model.Session) error {
	logger := logging.FromContext(ctx)

	logger.Debugw("removing sessoin", "userID", session.UserID, "deviceID", session.DeviceID, "serverID", session.ServerID)

	s.db.SessionSet(ctx, session)
	return nil
}

func (s *SessionService) IsUserActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	logger := logging.FromContext(ctx)

	logger.Debugw("is active", "userID", userID)

	sessions, err := s.db.SessionsGetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	if len(sessions) > 0 {
		return true, nil
	}

	return false, nil
}
