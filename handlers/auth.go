package handlers

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/Schtolc/mooncore/models"
	"encoding/json"
	"fmt"
)

type JwtClaims struct {
	Name string    `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}
var SigningKey = []byte("secret")

func Register(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
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
		// insert into db
		fmt.Println(dbUser)
		if dbc := db.Create(dbUser); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return c.JSON(http.StatusInternalServerError, InternalError)
		}
		// return response
		return c.JSON(http.StatusOK, &Resp{
			Code:    "200",
			Message: "You are registered. Welcome "+dbUser.Name,
		})
	}
}
func Login(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userAttr := models.User{}
		if err := json.NewDecoder(c.Request().Body).Decode(&userAttr); err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusInternalServerError, InternalError)
		}
		dbUser := &models.User{}
		db.Where("name = ? AND password = ?", userAttr.Name, userAttr.Password).First(dbUser)
		fmt.Println(dbUser)
		if  dbUser == nil {
			logrus.Info("login unregistered user")
			return c.JSON(http.StatusBadRequest, NeedRegistration)
		}
		tokenString, err := CreateJwtToken(dbUser)
		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusInternalServerError, InternalError)
		}
		// return response
		return c.JSON(http.StatusOK, &Resp{
			Code:    "200",
			Message: tokenString,
		})
	}
}
func Restricted(c echo.Context) error {
	q := c.Get("user")
	fmt.Println(q)
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func CheckJwtToken (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("in checker")
		auth := c.Request().Header.Get("Authorization")[7:]
		fmt.Println(auth)
		token, err := jwt.ParseWithClaims(auth, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return SigningKey, nil
		})
		if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
			fmt.Printf("%v %v %v", claims.Name, claims.Password, claims.StandardClaims.ExpiresAt)
		} else {
			fmt.Println(err)
		}
		return next(c)
	}
}

func CreateJwtToken(user *models.User) (tokenString string, err error) {
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