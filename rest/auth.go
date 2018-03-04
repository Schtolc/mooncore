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
		userData := c.Get(utils.UserKey)
		user := &models.User{}
		var err error
		if userData != nil {
			email := userData.(*jwt.Token).Claims.(*models.JwtClaims).Email
			user, err = dao.GetUserByEmail(email)
			if err != nil || user == nil {
				logrus.Info("User was not found in the database when checking token: ", email)
				return utils.SendResponse(c, http.StatusBadRequest, "Bad token")
			}
		} else {
			user = models.AnonUser
		}
		c.Set(utils.UserKey, user)
		return next(c)
	}
}
