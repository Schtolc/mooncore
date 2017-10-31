package handlers

import (
	"encoding/json"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type jwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

var signingKey = []byte(dependencies.ConfigInstance().Server.Auth.Secret)

// GetJwtConfig return configuration for jwt registration
func GetJwtConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		SigningMethod: "HS256",
		Claims:        &jwtClaims{},
		SigningKey:    signingKey,
	}
}

// SignUp registers new users
func SignUp(c echo.Context) error {
	userAttr := &models.UserAuth{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
		logrus.Error(c.Request().Body, err)
		return sendResponse(c, http.StatusBadRequest, "Bad json")
	}

	if err := userAttr.Validate(true); err != nil {
		logrus.Info(err)
		return sendResponse(c, http.StatusBadRequest, err.Error())
	}

	// check body with wrong fields
	user := &models.UserAuth{}
	dependencies.DBInstance().Where("name = ?", userAttr.Name).First(user)
	if *user != (models.UserAuth{}) {
		logrus.Info("User already exists: ", userAttr.Name)
		return sendResponse(c, http.StatusBadRequest, "User already exists: "+userAttr.Name)
	}

	if dbc := dependencies.DBInstance().Create(&models.UserAuth{
		Name:     userAttr.Name,
		Email:    userAttr.Email,
		Password: userAttr.Password,
	}); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return internalServerError(c)
	}

	return sendResponse(c, http.StatusOK, "")
}

// SignIn users; return auth token
func SignIn(c echo.Context) error {
	userAttr := models.UserAuth{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
		logrus.Error(err)
		return sendResponse(c, http.StatusBadRequest, "Bad json")
	}

	if err := userAttr.Validate(false); err != nil {
		logrus.Info(err)
		return sendResponse(c, http.StatusBadRequest, err.Error())
	}

	dbUser := &models.UserAuth{}
	dbc := dependencies.DBInstance().Where("name = ? AND password = ?", userAttr.Name, userAttr.Password).First(dbUser)
	if dbc.Error != nil {
		logrus.Info("Unregistered user: ", userAttr.Name)
		return sendResponse(c, http.StatusBadRequest, "Unregistered user")
	}

	tokenString, err := createJwtToken(dbUser)
	if err != nil {
		logrus.Error(err)
		return internalServerError(c)
	}

	return sendResponse(c, http.StatusOK, tokenString)
}

func createJwtToken(user *models.UserAuth) (tokenString string, err error) {
	claims := &jwtClaims{
		user.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(dependencies.ConfigInstance().Server.Auth.Lifetime)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if tokenString, err = token.SignedString(signingKey); err != nil {
		return "", err
	}
	return tokenString, nil
}
