package service

import (
	"auth-service/internal/model"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func New(secret []byte) *Service {
	return &Service{secretKey: secret}
}

type Service struct {
	secretKey []byte
}

func (s *Service) NewTokens(id string, role string) (*model.Tokens, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"iss": id, "iat": time.Now().Unix(), "exp": (time.Now().Add(time.Hour * 24 * 7)).Unix(), "type": "refresh", "role": role})
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"iss": id, "iat": time.Now().Unix(), "exp": (time.Now().Add(time.Hour * 3)).Unix(), "type": "access", "role": role})
	accessString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}
	refreshString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}
	tokens := &model.Tokens{
		Access:  accessString,
		Refresh: refreshString,
	}
	return tokens, nil
}

func (s *Service) Verify(token string) (int, string, error) {
	accessToken, err := jwt.ParseWithClaims(token, &model.Token{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return -1, "", err
	}
	if claims, ok := accessToken.Claims.(*model.Token); ok {
		iss, err := claims.GetIssuer()
		role := claims.Role
		if err != nil {
			return -1, "", err
		}
		id, err := strconv.Atoi(iss)
		if err != nil {
			return -1, "", err
		}
		return id, role, nil
	}
	return -1, "", err
}
