package handlers

import (
	"encoding/json"
	"github.com/Schtolc/mooncore/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// SigningKey use for generation token
var SigningKey = []byte("secret")

// JwtClaims - custom config for jwt
type JwtClaims struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GetJwtConfig return configuration for jwt registration
func GetJwtConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		SigningMethod: "HS256",
		Claims:        &JwtClaims{},
		SigningKey:    SigningKey,
	}
}

// SignUp registers new users
func (h *Handler) SignUp(c echo.Context) error {
	userAttr := &models.User{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, internalError)
	}

	user := &models.User{}
	h.DB.Where("name = ? AND password = ?", userAttr.Name, userAttr.Password).First(user)
	if *user != (models.User{}) {
		logrus.WithFields(logrus.Fields{
			"Name":     userAttr.Name,
			"Password": userAttr.Password,
		}).Info("User already exists")
		return c.JSON(http.StatusBadRequest, userAlreadyExists)
	}

	if dbc := h.DB.Create(&models.User{
		Name:     userAttr.Name,
		Email:    userAttr.Email,
		Password: userAttr.Password,
	}); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return c.JSON(http.StatusInternalServerError, internalError)
	}

	return c.NoContent(http.StatusOK)
}

// SignIn users; return auth token
func (h *Handler) SignIn(c echo.Context) error {
	userAttr := models.User{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, internalError)
	}
	dbUser := &models.User{}
	h.DB.Where("name = ? AND password = ?", userAttr.Name, userAttr.Password).First(dbUser)
	if *dbUser == (models.User{}) {
		logrus.WithFields(logrus.Fields{
			"Name":     userAttr.Name,
			"Password": userAttr.Password,
		}).Info("Unregistered user")
		return c.JSON(http.StatusBadRequest, needRegistration)
	}

	tokenString, err := createJwtToken(dbUser)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, internalError)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"Code": "200",
		"Token": tokenString,
	})
}

// SignOut remove user from our database
func (h *Handler) SignOut(c echo.Context) error {
	userAttr := models.User{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, internalError)
	}
	h.DB.Where("name = ? AND password = ?", userAttr.Name, userAttr.Password).Delete(&models.User{})
	return c.NoContent(http.StatusOK)
}

// CheckJwtToken verifies the validity of token
func (h *Handler) CheckJwtToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")[7:]
		token, err := jwt.ParseWithClaims(auth, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return SigningKey, nil
		})
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusInternalServerError, internalError)
		}
		if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
			logrus.WithFields(logrus.Fields{
				"Name":      claims.Name,
				"ExpiresAt": claims.StandardClaims.ExpiresAt,
			}).Info("Token verification")

			dbUser := &models.User{}
			if h.DB.Where("name = ? AND password = ?", claims.Name, claims.Password).First(dbUser); *dbUser == (models.User{}) {
				logrus.WithFields(logrus.Fields{
					"Name":     claims.Name,
					"Password": claims.Password,
				}).Info("User was not found in the database when checking token")
				return c.JSON(http.StatusBadRequest, needRegistration)
			}
			return next(c)
		}
		logrus.Warn("CheckJwtToken failed. Token is invalid: %v", token)
		return c.JSON(http.StatusBadRequest, invalidToken)
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
