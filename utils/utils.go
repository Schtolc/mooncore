package utils

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// UserKey is a context key for current logged user object
const UserKey = "user"

// HashPassword hashes password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks password and hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var signingKey = []byte(dependencies.ConfigInstance().Server.Auth.Secret)

// GetJwtConfig return configuration for jwt registration
func GetJwtConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		SigningMethod: "HS256",
		Claims:        &models.JwtClaims{},
		SigningKey:    signingKey,
		Skipper: func(c echo.Context) bool {
			return len(c.Request().Header.Get(echo.HeaderAuthorization)) == 0
		},
	}
}

// CreateJwtToken creates a new token
func CreateJwtToken(user *models.User) (tokenString string, err error) {
	claims := &models.JwtClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(dependencies.ConfigInstance().Server.Auth.Lifetime)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if tokenString, err = token.SignedString(signingKey); err != nil {
		return "", err
	}
	return tokenString, nil
}
