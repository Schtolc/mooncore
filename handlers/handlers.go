package handlers

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Resp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Ping server
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, &Resp{
		Code:    "200",
		Message: "ECHO_PING",
	})
}

// Ping database: create metric un
func PingDb(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		m := &models.Metric{
			Path: c.Path(),
			Time: gorm.NowFunc(),
		}
		if dbc := db.Create(m); dbc.Error != nil {
			log.Error(dbc.Error)
			return c.JSON(http.StatusInternalServerError, &Resp{
				Code:    "500",
				Message: "InternalError",
			})
		}
		return c.JSON(http.StatusOK, &Resp{
			Code:    "200",
			Message: strconv.Itoa(m.Id),
		})
	}
}
