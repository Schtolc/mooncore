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
	"github.com/Schtolc/mooncore/dependencies"
)


type jwtClaims struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var signingKey = []byte("secret")

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
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, internalError)
	}
	// check body with wrong fields
	user := &models.UserAuth{}
	dependencies.DBInstance().Where("name = ? AND password = ?", userAttr.Name, userAttr.Password).First(user)
	if *user != (models.UserAuth{}) {
		logrus.Info("User already exists: %s", userAttr.Name)
		return c.JSON(http.StatusBadRequest, userAlreadyExists)
	}

	if dbc := dependencies.DBInstance().Create(&models.UserAuth{
		Name:     userAttr.Name,
		Email:    userAttr.Email,
		Password: userAttr.Password,
	}); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return c.JSON(http.StatusInternalServerError, internalError)
	}

	return c.JSON(http.StatusOK, &Response{
		Code: OK,
	})
}

// SignIn users; return auth token
func SignIn(c echo.Context) error {
	userAttr := models.UserAuth{}
	if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, internalError)
	}
	dbUser := &models.UserAuth{}
	dependencies.DBInstance().Where("name = ? AND password = ?", userAttr.Name, userAttr.Password).First(dbUser)
	if *dbUser == (models.UserAuth{}) {
		logrus.Info("Unregistered user: %s", userAttr.Name)
		return c.JSON(http.StatusBadRequest, needRegistration)
	}

	tokenString, err := createJwtToken(dbUser)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, internalError)
	}

	return c.JSON(http.StatusOK, &Response{
		Code: OK,
		Body: tokenString,
	})
}

// CheckJwtToken verifies the validity of token
func CheckJwtToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")[7:]
		token, err := jwt.ParseWithClaims(auth, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusInternalServerError, internalError)
		}
		if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
			logrus.Info("Token verification: %s %s", claims.Name, claims.StandardClaims.ExpiresAt)

			dbUser := &models.UserAuth{}
			dependencies.DBInstance().Where("name = ? AND password = ?", claims.Name, claims.Password).First(dbUser)
			if *dbUser == (models.UserAuth{}) {
				logrus.Info("User was not found in the database when checking token: %s", claims.Name)
				return c.JSON(http.StatusBadRequest, needRegistration)
			}
			return next(c)
		}
		logrus.Warn("CheckJwtToken failed. Token is invalid: %v", token)
		return c.JSON(http.StatusBadRequest, invalidToken)
	}
}

func createJwtToken(user *models.UserAuth) (tokenString string, err error) {
	claims := &jwtClaims{
		user.Name,
		user.Password,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if tokenString, err = token.SignedString(signingKey); err != nil {
		return "", err
	}
	return tokenString, nil
}
