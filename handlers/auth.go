package handlers

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"time"

	"encoding/json"
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// JwtClaims - custom config for jwt
type JwtClaims struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// SigningKey use for generation token
var SigningKey = []byte("secret")

// Register method for new users
func (h *Handler) Register(c echo.Context) error {
	userAttr := models.User{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, InternalError)
	}
	dbUser := &models.User{
		Name:     userAttr.Name,
		Email:    userAttr.Email,
		Password: userAttr.Password,
	}

	if dbc := h.DB.Create(dbUser); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return c.JSON(http.StatusInternalServerError, InternalError)
	}

	return c.JSON(http.StatusOK, &Resp{
		Code:    "200",
		Message: "You are registered. Welcome " + dbUser.Name,
	})
}

// Login method give token to register user
func (h *Handler) Login(c echo.Context) error {
	userAttr := models.User{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, InternalError)
	}
	dbUser := &models.User{}
	h.DB.Where("name = ? AND password = ?", userAttr.Name, userAttr.Password).First(dbUser)

	if dbUser.IsEmpty() {
		logrus.WithFields(logrus.Fields{
			"Name":     userAttr.Name,
			"Password": userAttr.Password,
		}).Info("Unregistered user")
		return c.JSON(http.StatusBadRequest, NeedRegistration)
	}
	tokenString, err := createJwtToken(dbUser)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, InternalError)
	}

	return c.JSON(http.StatusOK, &Resp{
		Code:    "200",
		Message: tokenString,
	})
}

// CheckJwtToken for validation and existing in db
func (h *Handler) CheckJwtToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")[7:]

		token, err := jwt.ParseWithClaims(auth, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return SigningKey, nil
		})
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusInternalServerError, InternalError)
		}
		if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
			logrus.WithFields(logrus.Fields{
				"Name":      claims.Name,
				"ExpiresAt": claims.StandardClaims.ExpiresAt,
			}).Info("check jwt token")

			dbUser := &models.User{}
			h.DB.Where("name = ? AND password = ?", claims.Name, claims.Password).First(dbUser)
			if dbUser.IsEmpty() {
				logrus.WithFields(logrus.Fields{
					"Name":     claims.Name,
					"Password": claims.Password,
				}).Info("Unregistered user")
				return c.JSON(http.StatusBadRequest, NeedRegistration)
			}
			return next(c)
		} else {
			logrus.Warn("CheckJwtToken failed. Token is invalid: %v", token)
			return c.JSON(http.StatusBadRequest, InvalidToken)
		}
	}
}

func createJwtToken(user *models.User) (tokenString string, err error) {
	claims := JwtClaims{
		user.Name,
		user.Password,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if tokenString, err = token.SignedString(SigningKey); err != nil {
		return "", err
	}
	return tokenString, nil
}
