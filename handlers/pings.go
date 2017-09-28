package handlers

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"github.com/Schtolc/mooncore/dependencies"
)


// Ping is a simple handler for checking if server is up and running.
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, &Response{
		Code:  OK,
		Body: "ECHO_PING",
	})
}

func PingAuth (c echo.Context) error {
	return c.JSON(http.StatusOK, &Response{
		Code:  OK,
		Body: "ECHO_AUTH_PING",
	})
}

// PingDb is a simple handler for checking if database is up and running.
func PingDb(c echo.Context) error {
	m := &models.Mock{
		Path: c.Path(),
		Time: gorm.NowFunc(),
	}
	if dbc := dependencies.DBInstance().Create(m); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return c.JSON(http.StatusInternalServerError, internalError)
	}
	return c.JSON(http.StatusOK, &Response{
		Code: OK,
		Body: strconv.Itoa(m.ID),
	})
}
