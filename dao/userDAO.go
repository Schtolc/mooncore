package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// SignIn checks user credentials and returns token
func SignIn(email, password string) (*models.Token, error) {
	user := &models.User{}

	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		logrus.Info("Unregistered user: ", email)
		return nil, nil
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		logrus.Info("Unregistered user: ", email)
		return nil, nil
	}

	tokenString, err := utils.CreateJwtToken(user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &models.Token{Token: tokenString}, nil
}

// GetUserByID returns user by id
func GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}

	if err := db.First(&user, id).Error; err != nil {
		logrus.Info("User not found: ", id)
		return nil, err
	}

	return user, nil
}

// GetUserByEmail returns user by email
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := db.Where("email = ? ", email).First(user).Error; err != nil {
		logrus.Info("User not found: ", email)
		return nil, err
	}

	return user, nil
}

func createUser(email, password string, tx *gorm.DB) (*models.User, error) {
	passwordHash, err := utils.HashPassword(password)

	if err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         0,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}

	return user, nil
}

func deleteUser(id int64) error {
	if err := db.Delete(models.User{ID: id}).Error; err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
