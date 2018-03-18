package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// GetClientByID returns client by id
func GetClientByID(id int64) (*models.Client, error) {
	user := &models.User{}
	client := &models.Client{}
	if err := db.First(client, id).Error; err != nil {
		logrus.Error(err)
		return nil, nil
	}
	if err := db.First(user, client.UserID).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	client.User = *user
	return client, nil
}

// CreateClient creates new Client
func CreateClient(email, passwordHash string) (*models.Client, error) {
	tx := db.Begin()
	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         models.ClientRole,
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}

	client := &models.Client{
		UserID: user.ID,
	}
	if err := tx.Create(client).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}
	tx.Commit()
	return client, nil
}

// DeleteClient deletes client
func DeleteClient(id int64) error {
	client, err := GetClientByID(id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	tx := db.Begin()
	if err = tx.Delete(models.Client{ID: client.ID}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}

	if err := tx.Delete(models.User{ID: client.UserID}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}
	tx.Commit()
	return nil
}
