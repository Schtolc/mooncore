package handlers

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// Resp represents simple json response with code and message.
type Resp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Handler type: params for each handler
type Handler struct {
	DB *gorm.DB
}

// Init Handler with parameters : db
func Init(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

// Ping is a simple handler for checking if server is up and running.
func (h *Handler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "ECHO_PING")
}

// PingDb is a simple handler for checking if database is up and running.
func (h *Handler) PingDb(c echo.Context) error {
	m := &models.Mock{
		Path: c.Path(),
		Time: gorm.NowFunc(),
	}
	if dbc := h.DB.Create(m); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return c.JSON(http.StatusInternalServerError, internalError)
	}
	return c.JSON(http.StatusOK, "Database pinged: " + strconv.Itoa(m.ID))
}
