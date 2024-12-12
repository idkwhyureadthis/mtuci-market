package service

import (
	"auth-service/internal/model"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

func (s *Service) NewTokens(id string) (*model.Tokens, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"iss": id, "iat": time.Now().Unix(), "exp": (time.Now().Add(time.Hour * 6)).Unix(), "type": "refresh"})
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"iss": id, "iat": time.Now().Unix(), "exp": (time.Now().Add(time.Second * 1)).Unix(), "type": "access"})
	accessString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}
	refreshString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}
	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(refreshString))
	cryptedRefresh := hex.EncodeToString(h.Sum(nil))
	tokens := &model.Tokens{
		Access:         accessString,
		Refresh:        refreshString,
		CryptedRefresh: cryptedRefresh,
	}
	return tokens, nil
}

func (s *Service) Verify(token string) (int, error) {
	accessToken, err := jwt.ParseWithClaims(token, &model.Token{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return -1, err
	}
	if claims, ok := accessToken.Claims.(*model.Token); ok {
		iss, err := claims.GetIssuer()
		if err != nil {
			return -1, err
		}
		id, err := strconv.Atoi(iss)
		if err != nil {
			return -1, err
		}
		return id, nil
	}
	return -1, err
}
