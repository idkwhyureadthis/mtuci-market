package endpoint

import (
	"api-service/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Endpoint struct {
	s *service.Service
	c *echo.Echo
}

func New() *Endpoint {
	e := Endpoint{
		c: echo.New(),
		s: service.New(),
	}
	e.c.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{"POST", "GET"},
		MaxAge:       600,
	}))
	e.setupHandlers()
	return &e
}

func (e *Endpoint) Start(port string) {
	e.c.Logger.Fatal(e.c.Start(port))
}

func (e *Endpoint) setupHandlers() {
	e.c.POST("/signup", e.signUp)
	e.c.POST("/login", e.logIn)
	e.c.POST("/verify", e.verify)
	e.c.POST("/create_card", e.createCard)
	e.c.GET("/get_cards", e.getCards)
	e.c.GET("/products_of", e.productsOf)
	e.c.POST("/delete_card/:id", e.deleteCard)
	e.c.POST("/create_moderator", e.createModerator)
	e.c.GET("/on_moderation", e.onModeration)
	e.c.POST("/accept/:id", e.accept)
	e.c.POST("/reject/:id", e.reject)
}

func (e *Endpoint) accept(c echo.Context) error {
	id := c.Param("id")
	resp, err := http.Post("http://db-service:8081/accept/"+id, "", nil)
	if err != nil || resp.StatusCode != 200 {
		return c.String(http.StatusInternalServerError, "failed to change card status")
	}
	return c.String(http.StatusOK, "card successfully updated")
}

func (e *Endpoint) reject(c echo.Context) error {
	id := c.Param("id")
	resp, err := http.Post("http://db-service:8081/reject/"+id, "", nil)
	if err != nil || resp.StatusCode != 200 {
		return c.String(http.StatusInternalServerError, "failed to change card status")
	}
	return c.String(http.StatusOK, "card successfully updated")
}

func (e *Endpoint) onModeration(c echo.Context) error {
	card, err := e.s.OnModeration()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, card)
}

func (e *Endpoint) createModerator(c echo.Context) error {
	data := struct {
		Name     string `json:"name"`
		Telegram string `json:"telegram"`
		Password string `json:"password"`
		Login    string `json:"login"`
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = e.s.CreateModerator(data)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "moderator created")
}

func (e *Endpoint) deleteCard(c echo.Context) error {
	id := c.Param("id")
	err := e.s.DeleteCard(id)
	if err != nil {
		c.Logger().Printf("failed to delete card")
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "card deleted successfully")
}

func (e *Endpoint) productsOf(c echo.Context) error {
	idString := c.QueryParam("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	cards, err := e.s.ProductsOf(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"cards": cards})
}

func (e *Endpoint) getCards(c echo.Context) error {
	cards, err := e.s.GetCards()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"cards": cards})
}

func (e *Endpoint) createCard(c echo.Context) error {
	form, err := c.MultipartForm()
	if form == nil || err != nil {
		c.Logger().Print("no data provided")
		return c.String(http.StatusBadRequest, "no data provided")
	}
	newTokens, err := e.s.CreateCard(form)
	if err != nil && strings.HasSuffix(err.Error(), "authorized") {
		return c.String(http.StatusUnauthorized, "user not authorized")
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"tokens": newTokens})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"tokens": newTokens})
}

func (e *Endpoint) verify(c echo.Context) error {
	data := struct {
		Refresh string `json:"refresh"`
		Access  string `json:"access"`
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	userData, _, err := e.s.Verify(data.Access, data.Refresh)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userData)
}

func (e *Endpoint) signUp(c echo.Context) error {
	data := struct {
		Login      string `json:"login"`
		Password   string `json:"password"`
		DormNumber string `json:"dorm_number"`
		Room       string `json:"room"`
		Name       string `json:"name"`
		Telegram   string `json:"telegram"`
	}{}

	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	tokens, err := e.s.SignUp(data.Login, data.Password, data.DormNumber, data.Room, data.Name, data.Telegram)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, tokens)
}

func (e *Endpoint) logIn(c echo.Context) error {
	data := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&data)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	tokens, err := e.s.LogIn(data.Login, data.Password)
	if err != nil {
		c.Logger().Print(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(200, tokens)
}
