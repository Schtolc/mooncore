package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// GetClientByID returns client by id
func GetClientByID(id int64) (*models.Client, error) {
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

// CreateClient creates new client
func CreateClient(username, email, password, name string, photoID int64) (*models.Client, error) {
	tx := db.Begin()

	user, err := createUser(email, password, tx)
	if err != nil {
		return nil, err
	}

	client := &models.Client{
		UserID:  user.ID,
		Name:    name,
		PhotoID: photoID,
	}

	if err := tx.Create(client).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return client, nil
}

// DeleteClient deletes client
func DeleteClient(id int64) error {
	client, err := GetClientByID(id)

	if err != nil {
		return err
	}

	err = db.Delete(models.Client{ID: id}).Error

	if err != nil {
		return err
	}

	return deleteUser(client.UserID)
}
