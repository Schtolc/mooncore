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

// CreateClient creates new client
func CreateClient(username, email, password, name string, photoID int64) (*models.Client, error) {
	tx := db.Begin()

	user, err := createUser(email, password, models.Roles["Client"], tx)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	client := &models.Client{
		UserID:  user.ID,
		Name:    name,
		PhotoID: photoID,
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

	err = db.Delete(models.Client{ID: id}).Error
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = deleteUser(client.UserID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
