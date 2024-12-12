package handler

import (
	"auth-service/pkg/db"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *db.DB
	e  *echo.Echo
}

func New(connURL string) *Handler {
	e := echo.New()
	h := &Handler{db: db.New(connURL), e: e}
	h.setupHandlers(h.e)
	return h
}

func (h *Handler) setupHandlers(e *echo.Echo) {
	e.GET("/newuser", h.createNewUser)
}

func (h *Handler) createNewUser(c echo.Context) error {
	return nil
}
