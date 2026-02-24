package jwt

import (
	"bluebell/settings"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

var mySecret = []byte("godblf")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenJwtToken(userId int64, username string) (string, error) {
	c := &MyClaims{
		UserID:   userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(settings.GlobalConfig.AuthConfig.JwtExpire) * time.Hour).Unix(), // 过期时间
			Issuer: "bluebell", // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(mySecret)
}

type AuthError struct {
	name string
}

func NewAuthError(name string) *AuthError {
	return &AuthError{name: name}
}

func (a *AuthError) Error() string {
	return a.name
}

var (
	AuthErrorInvalid = NewAuthError("invalid token")
)

func ParseJwtToken(token string) (*MyClaims, error) {
	mc := &MyClaims{}
	claims, err := jwt.ParseWithClaims(token, mc, func(tk *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		zap.L().Error("jwt.ParseWithClaims failed", zap.Error(err))
		return nil, fmt.Errorf("jwt.ParseWithClaims failed, err: %v", err)
	}
	if !claims.Valid {
		zap.L().Info("invalid token")
		return nil, AuthErrorInvalid
	}
	return mc, nil

}
