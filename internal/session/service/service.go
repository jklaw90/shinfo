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
	cfg config.Provider
}

func New(ctx context.Context, cfg config.Provider) (*SessionService, error) {
	logger := logging.FromContext(ctx)

	logger.Debugw("session service initializing")

	db, err := database.New(ctx, cfg)
	if err != nil {
		logger.Debugw("session service initialization failed")
		return nil, err
	}

	return &SessionService{
		db:  db,
		cfg: cfg,
	}, nil
}

func (s *SessionService) Register(ctx context.Context, session model.Session) error {
	logger := logging.FromContext(ctx)

	logger.Debugw("registering session", "userID", session.UserID, "deviceID", session.DeviceID, "serverID", session.ServerID)

	if err := s.db.SessionSet(ctx, session); err != nil {
		return err
	}

	return nil
}

func (s *SessionService) Deregister(ctx context.Context, session model.Session) error {
	logger := logging.FromContext(ctx)

	logger.Debugw("deregistering session", "userID", session.UserID, "deviceID", session.DeviceID, "serverID", session.ServerID)

	if err := s.db.SessionSet(ctx, session); err != nil {
		return err
	}

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
