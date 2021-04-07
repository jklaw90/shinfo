package auth

import (
	"context"

	"github.com/jklaw90/shinfo/pkg/model"
)

type Service interface {
	MakeToken(ctx context.Context, user model.User) (string, error)
	VerifyToken(ctx context.Context, jwt []byte) (model.User, error)
}
