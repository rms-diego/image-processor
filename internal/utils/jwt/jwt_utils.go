package jwtutils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rms-diego/image-processor/internal/validations"
	"github.com/rms-diego/image-processor/pkg/config"
)

type JwtUtilsInterface interface {
	GenerateToken(user validations.UserFound) (*string, error)
	ValidateAndDecodeToken(token string) (*JwtDecoded, error)
}

type jwtUtils struct{}

type JwtDecoded struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

func NewJwtUtils() JwtUtilsInterface {
	return &jwtUtils{}
}

func (j *jwtUtils) GenerateToken(user validations.UserFound) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.Env.JWT_SECRET))
	if err != nil {
		return nil, err
	}

	return &tokenStr, nil
}

func (j *jwtUtils) ValidateAndDecodeToken(tokenStr string) (*JwtDecoded, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}

		return []byte(config.Env.JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenExpired
	}

	t := JwtDecoded{
		ID:       token.Claims.(jwt.MapClaims)["id"].(string),
		Username: token.Claims.(jwt.MapClaims)["username"].(string),
		Exp:      int64(token.Claims.(jwt.MapClaims)["exp"].(float64)),
	}

	return &t, nil
}
