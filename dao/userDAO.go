package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
)

func GetMasterById(id int64) (*models.Master, error) {
	user := &models.User{}
	master := &models.Master{}
	if dbc := db.First(master, id); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, nil
	}
	if dbc := db.First(user, master.UserID); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	master.User = *user
	return master, nil
}

func GetClientById(id int64) (*models.Client, error) {
	user := &models.User{}
	client := &models.Client{}
	if dbc := db.First(client, id); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, nil
	}
	if dbc := db.First(user, client.UserID); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	client.User = *user
	return client, nil
}

func CreateMaster(username, email, password, name string, addressId, photoId int64) (*models.Master, error) {
	tx := db.Begin()

	passwordHash, err := HashPassword(password)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         0,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	master := &models.Master{
		UserID:    user.ID,
		Name:      name,
		AddressID: addressId,
		PhotoID:   photoId,
	}

	if err := tx.Create(master).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return master, nil
}

func CreateClient(username, email, password, name string, photoId int64) (*models.Client, error) {
	tx := db.Begin()

	passwordHash, err := HashPassword(password)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         0,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	client := &models.Client{
		UserID:  user.ID,
		Name:    name,
		PhotoID: photoId,
	}

	if err := tx.Create(client).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return client, nil
}

func SignIn(username, email, password string) (*models.Token, error) {
	user := &models.User{}

	if err := db.Where("username = ? AND email =? AND password = ?", username, email, password).First(user).Error; err != nil {
		logrus.Info("Unregistered user: ", username)
		return nil, nil
	}

	tokenString, err := CreateJwtToken(user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &models.Token{Token: tokenString}, nil
}

// TODO remove all to the end of the file

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

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

func CreateJwtToken(user *models.User) (tokenString string, err error) {
	claims := &models.JwtClaims{
		Name: user.Username,
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
