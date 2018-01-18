package rest

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// Ping is a simple handler for checking if server is up and running.
func Ping(c echo.Context) error {
	return utils.SendResponse(c, http.StatusOK, "ECHO_PING")
}

// PingAuth is a handler for checking if authorization works.
func PingAuth(c echo.Context) error {
	user := c.Get("user").(*models.User)
	return utils.SendResponse(c, http.StatusOK, user.Username)
}

// PingDb is a simple handler for checking if database is up and running.
func PingDb(c echo.Context) error {
	m := &models.Mock{
		Path: c.Path(),
		Time: gorm.NowFunc(),
	}
	if err := dependencies.DBInstance().Create(m).Error; err != nil {
		logrus.Error(err)
		return utils.InternalServerError(c)
	}
	return utils.SendResponse(c, http.StatusOK, strconv.Itoa(m.ID))
}
