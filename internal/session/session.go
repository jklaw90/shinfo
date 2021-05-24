package session

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, userID, deviceID uuid.UUID) error
	Deregister(ctx context.Context, userID, deviceID uuid.UUID) error
	IsUserActive(ctx context.Context, userID uuid.UUID) (bool, error)
}
