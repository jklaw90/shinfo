package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"github.com/jklaw90/shinfo/pkg/model"
)

type JWT struct {
	UserID  gocql.UUID `json:"uid"`
	Name    string     `json:"name,omitempty"`
	Picture string     `json:"picture,omitempty"`
	jwt.StandardClaims
}

func (j *JWT) GetUser() model.User {
	return model.User{
		ID:     j.UserID,
		Name:   j.Name,
		Avatar: j.Picture,
	}
}

func NewClaims(userID gocql.UUID, name, picture string, duration time.Duration) jwt.Claims {
	claims := JWT{
		UserID:  userID,
		Name:    name,
		Picture: picture,
	}
	claims.Issuer = "shinfo.io"
	claims.IssuedAt = time.Now().Unix()
	claims.NotBefore = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(duration).Unix()
	return claims
}

func MintJWT(secret []byte, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

func Parse(secret []byte, tokenString string) (JWT, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing alg: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return JWT{}, err
	}

	claims, ok := token.Claims.(*JWT)
	if !ok || token.Valid {
		return JWT{}, errors.New("invalid claims")
	}

	return *claims, err
}
