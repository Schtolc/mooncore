package handlers

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// Ping is a simple handler for checking if server is up and running.
func Ping(c echo.Context) error {
	return sendResponse(c, http.StatusOK, "ECHO_PING")
}

// PingAuth is a handler for checking if authorization works.
func PingAuth(c echo.Context) error {
	user := c.Get("user").(*models.UserAuth)
	return sendResponse(c, http.StatusOK, user.Name)
}

// PingDb is a simple handler for checking if database is up and running.
func PingDb(c echo.Context) error {
	m := &models.Mock{
		Path: c.Path(),
		Time: gorm.NowFunc(),
	}
	if dbc := dependencies.DBInstance().Create(m); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return internalServerError(c)
	}
	return sendResponse(c, http.StatusOK, strconv.Itoa(m.ID))
}
