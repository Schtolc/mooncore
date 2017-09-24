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

// Ping is a simple handler for checking if server is up and running.
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, &Resp{
		Code:    "200",
		Message: "ECHO_PING",
	})
}

// PingDb is a simple handler for checking if database is up and running.
func PingDb(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		m := &models.Mock{
			Path: c.Path(),
			Time: gorm.NowFunc(),
		}
		if dbc := db.Create(m); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return c.JSON(http.StatusInternalServerError, &Resp{
				Code:    "500",
				Message: "InternalError",
			})
		}
		return c.JSON(http.StatusOK, &Resp{
			Code:    "200",
			Message: strconv.Itoa(m.ID),
		})
	}
}
