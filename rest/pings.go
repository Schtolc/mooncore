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

// Ping is a simple handler for checking server is up and running
func Ping(c echo.Context) error {
	return utils.SendResponse(c, http.StatusOK, "ECHO_PING")
}

// PingAuth is a handler for checking authorization works
func PingAuth(c echo.Context) error {
	user := c.Get(utils.UserKey).(*models.User)
	return utils.SendResponse(c, http.StatusOK, user.Email)
}

// PingDb is a simple handler for checking database is up and running
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
