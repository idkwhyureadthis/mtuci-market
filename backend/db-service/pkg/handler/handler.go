package handler

import (
	"auth-service/pkg/db"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *db.DB
	e  *echo.Echo
}

func New(connURL string) *Handler {
	e := echo.New()
	h := &Handler{db: db.New(connURL), e: e}
	h.createAdmin()
	h.setupHandlers(h.e)
	return h
}

func (h *Handler) createAdmin() {
	h.db.CreateNewUser(context.Background(), "admin", "", os.Getenv("ADMIN_PASSWORD"), "", "", os.Getenv("ADMIN_NAME"), "admin")
}

func (h *Handler) setupHandlers(e *echo.Echo) {
	e.POST("/create_user", h.createNewUser)
	e.POST("/send_message", h.sendMessage)
	e.POST("/delete_user", h.deleteUser)
	e.POST("/create_card", h.createCard)
	e.GET("/on_moderation", h.getOnModeration)
	e.GET("/moderated", h.getModerated)
	e.POST("/check_password", h.verifyPassword)
	e.POST("/create_moderator", h.createModerator)
	e.POST("/set_refresh", h.setRefresh)
	e.GET("/get_data", h.getUserData)
	e.GET("/verify_refresh", h.verifyRefresh)
	e.GET("/image", h.getImage)
	e.GET("/cards_of", h.cardsOf)
	e.POST("/delete_card/:id", h.deleteCard)
	e.POST("/reject/:id", h.reject)
	e.POST("/accept/:id", h.accept)
}

func (h *Handler) reject(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.db.ChangeCardStatus(c.Request().Context(), "rejected", id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "status change succesfully")
}

func (h *Handler) accept(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.db.ChangeCardStatus(c.Request().Context(), "accepted", id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "status change succesfully")
}

func (h *Handler) deleteCard(c echo.Context) error {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.db.DeleteCard(id)
	if err != nil {
		c.Logger().Printf("failed to delete card")
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "card deleted successfully")
}

func (h *Handler) Start(port string) {
	h.e.Logger.Fatal(h.e.Start(":" + port))
}

func (h *Handler) cardsOf(c echo.Context) error {
	idString := c.QueryParam("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	cards, err := h.db.CardsOf(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, cards)
}

func (h *Handler) getImage(c echo.Context) error {
	blob, err := h.db.GetImage(c.QueryParam("name"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"blob": blob,
	})
}

func (h *Handler) verifyRefresh(c echo.Context) error {
	idString, token := c.QueryParam("id"), c.QueryParam("refresh")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.String(http.StatusBadRequest, "wrong request data provided")
	}
	err = h.db.VerifyToken(id, token)
	if err != nil {
		return c.String(http.StatusNotAcceptable, "wrong token provided")
	}
	return c.String(http.StatusOK, "token correct")
}

func (h *Handler) setRefresh(c echo.Context) error {
	idString := c.QueryParam("id")
	refresh := c.QueryParam("refresh")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.db.SetRefreshKey(id, refresh)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "token set successfully")
}

func (h *Handler) getUserData(c echo.Context) error {
	idString := c.QueryParam("id")
	if idString == "" {
		c.Logger().Print("no id provided")
		return c.String(http.StatusBadRequest, "no id provided")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.Logger().Print("wrong data provided")
		return c.String(http.StatusBadRequest, "wrong data provided")
	}
	data, err := h.db.GetUserData(id)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, data)
}

func (h *Handler) createNewUser(c echo.Context) error {
	data := struct {
		Name       string `json:"name"`
		Telegram   string `json:"telegram"`
		Password   string `json:"password"`
		Room       string `json:"room"`
		DormNumber string `json:"dorm_number"`
		Login      string `json:"login"`
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	response := struct {
		Id int `json:"id"`
	}{}
	id, err := h.db.CreateNewUser(c.Request().Context(), data.Name, data.Telegram, data.Password, data.Room, data.DormNumber, data.Login, "user")
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	response.Id = id
	return c.JSON(200, response)
}

func (h *Handler) sendMessage(c echo.Context) error {
	data := struct {
		UserID  int    `json:"user_id"`
		To      int    `json:"send_to"`
		Content string `json:"content"`
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if data.UserID == data.To || data.Content == "" {
		c.Logger().Print("bad request data provided")
		return c.String(http.StatusBadRequest, "bad request data provided")
	}
	err = h.db.CreateMessage(c.Request().Context(), data.UserID, data.To, data.Content)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "message succesfully sent")
}

func (h *Handler) deleteUser(c echo.Context) error {
	data := struct {
		Id int `json:"id"`
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = h.db.DeleteUser(c.Request().Context(), data.Id)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "user successfully deleted")
}

func (h *Handler) createCard(c echo.Context) error {
	photos := make([]string, 0)
	form, err := c.MultipartForm()
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	idString := c.QueryParam("id")
	priceString := form.Value["price"]
	cardName := form.Value["name"]
	cardDescription := form.Value["description"]
	if len(priceString) != 1 || len(cardName) != 1 || len(cardDescription) != 1 {
		c.Logger().Print("wrong request data provided")
		return c.String(http.StatusBadRequest, "wrong request data provided")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, "wrong request data provided")
	}
	price, err := strconv.ParseFloat(priceString[0], 64)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	for _, file := range form.File["images"] {
		src, err := file.Open()
		if err != nil {
			c.Logger().Print(err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()
		bytes, err := io.ReadAll(src)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		ext := filepath.Ext(file.Filename)[1:]
		photoString := base64.StdEncoding.EncodeToString(bytes)
		photoString = fmt.Sprintf("data:image/%s; base64,", ext) + photoString
		photos = append(photos, photoString)
	}
	err = h.db.CreateCard(c.Request().Context(), id, cardName[0], price, photos, cardDescription[0])
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "card created successfully")
}

func (h *Handler) getOnModeration(c echo.Context) error {
	card, err := h.db.GetOnModeration()
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, card)
}

func (h *Handler) getModerated(c echo.Context) error {
	cards, err := h.db.GetModerated()
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, cards)
}

func (h *Handler) verifyPassword(c echo.Context) error {
	data := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	id, status, err := h.db.VerifyPassword(c.Request().Context(), data.Login, data.Password)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": status,
		"id":     id,
	})
}

func (h *Handler) createModerator(c echo.Context) error {
	data := struct {
		Name     string `json:"name"`
		Telegram string `json:"telegram"`
		Password string `json:"password"`
		Login    string `json:"login"`
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	response := struct {
		Id int `json:"id"`
	}{}
	id, err := h.db.CreateNewUser(c.Request().Context(), data.Name, data.Telegram, data.Password, "", "", data.Login, "moderator")
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	response.Id = id
	return c.JSON(200, response)
}

func (h *Handler) changeStatus(c echo.Context) error {
	idString := c.QueryParam("id")
	newStatus := c.QueryParam("status")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.db.ChangeCardStatus(c.Request().Context(), newStatus, id)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "status changed successfully")
}
