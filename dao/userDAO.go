package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/sirupsen/logrus"
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

	passwordHash, err := utils.HashPassword(password)

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

	passwordHash, err := utils.HashPassword(password)

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

	if err := db.Where("username = ? AND email = ?", username, email).First(user).Error; err != nil {
		logrus.Info("Unregistered user: ", username)
		return nil, nil
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		logrus.Info("Unregistered user: ", username)
		return nil, nil
	}

	tokenString, err := utils.CreateJwtToken(user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &models.Token{Token: tokenString}, nil
}

func Feed(offset, limit int64) ([]models.Master, error) {
	var masters []models.Master
	if err := db.Limit(limit).Offset(offset).Find(&masters).Error; err != nil {
		return nil, err
	}
	return masters, nil
}
