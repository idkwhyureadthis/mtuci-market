package endpoint

import (
	"auth-service/internal/model"
	"auth-service/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	SecretKey []byte
	echo      *echo.Echo
	s         *service.Service
}

func (e *Endpoint) Start(port string) {
	e.echo.Logger.Fatal(e.echo.Start(":" + port))
}

func SetupHandlers(e *Endpoint) {
	e.echo.GET("/new", e.NewTokens)
	e.echo.POST("/verify", e.VerifyToken)
	e.echo.POST("/generate", e.GenerateNewTokens)
}

func New(secret []byte) *Endpoint {
	e := echo.New()
	srv := &Endpoint{echo: e, SecretKey: secret, s: service.New(secret)}
	SetupHandlers(srv)
	return srv
}

func (e *Endpoint) NewTokens(c echo.Context) error {
	tokens, err := e.s.NewTokens(c.QueryParam("id"))
	if err != nil {
		c.Logger().Printf("failed to create new tokens: %v \n", err)
	}
	return c.JSON(http.StatusOK, tokens)
}

func (e *Endpoint) VerifyToken(c echo.Context) error {
	tokens := model.Tokens{}
	err := json.NewDecoder(c.Request().Body).Decode(&tokens)
	if err != nil {
		c.Logger().Printf("failed to parse request body : %v \n", err)
		return c.String(http.StatusBadRequest, "failed to parse request body")
	}
	id, err := e.s.Verify(tokens.Access)
	if errors.Is(err, jwt.ErrTokenExpired) {
		return c.String(http.StatusUnauthorized, "tokens expired")
	}
	if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrSignatureInvalid) {
		return c.String(http.StatusNotAcceptable, "token modified")
	}
	if err == nil {
		return c.JSON(http.StatusOK, model.VerificationResponse{Id: id})
	}
	return c.String(http.StatusBadRequest, err.Error())
}

func (s Endpoint) GenerateNewTokens(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "")
}
