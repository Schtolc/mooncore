package rest

import (
	"github.com/Schtolc/mooncore/dao"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Headers for option request
func Headers(c echo.Context) error {
	return utils.SendResponse(c, http.StatusOK, "")
}

// LoadUser is a middleware for load authorized user to context
func LoadUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get(utils.UserKey)
		if user == nil {
			return next(c)
		}
		email := user.(*jwt.Token).Claims.(*models.JwtClaims).Email

		user, err := dao.GetUserByEmail(email)
		if err != nil {
			logrus.Info("User was not found in the database when checking token: ", email)
			return utils.SendResponse(c, http.StatusBadRequest, "Bad token")
		}

		c.Set(utils.UserKey, user)
		return next(c)
	}
}
