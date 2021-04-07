package auth

import (
	"context"
	"time"

	"github.com/jklaw90/shinfo/pkg/jwt"
	"github.com/jklaw90/shinfo/pkg/model"
)

type AuthService struct {
	defaultDuration time.Duration
	secret          []byte
}

func (s *AuthService) MakeToken(ctx context.Context, user model.User) (string, error) {
	claims := jwt.NewClaims(user.ID, user.Name, user.Avatar, s.defaultDuration)
	token, err := jwt.MintJWT(s.secret, claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, token string) (model.User, error) {
	var users model.User

	claims, err := jwt.Parse(s.secret, token)
	if err != nil {
		return users, err
	}

	return claims.GetUser(), nil
}
