package rest

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

const UserKey = "user"

// Headers for option request
func Headers(c echo.Context) error {
	return utils.SendResponse(c, http.StatusOK, "")
}

// LoadUser is a middleware for load authorized user to context
func LoadUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(UserKey)
		if user == nil {
			return next(c)
		}
		username := user.(*jwt.Token).Claims.(*models.JwtClaims).Name
		dbUser := &models.User{}
		if dependencies.DBInstance().Where("name = ? ", username).First(dbUser).Error != nil {
			logrus.Info("User was not found in the database when checking token: ", username)
			return utils.SendResponse(c, http.StatusBadRequest, "Bad token")
		}
		c.Set(UserKey, dbUser)
		return next(c)
	}
}
