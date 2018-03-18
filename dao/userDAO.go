package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// GetUser from database
func GetUser(userParams *models.User) (*models.User, error) {
	user := &models.User{}
	if err := db.Where(userParams).First(user).Error; err != nil {
		logrus.Info("User not found")
		return nil, err
	}
	return user, nil
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
