package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
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

	user, err := createUser(email, password, tx)
	if err != nil {
		return nil, err
	}

	// TODO check addressID & photoId

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

	user, err := createUser(email, password, tx)
	if err != nil {
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

func DeleteMasterById(id int64) error {
	master, err := GetMasterById(id)

	if err != nil {
		return err
	}

	err = db.Delete(models.Master{ID: id}).Error

	if err != nil {
		return err
	}

	return deleteUserById(master.UserID)

}

func DeleteClientById(id int64) error {
	client, err := GetClientById(id)

	if err != nil {
		return err
	}

	err = db.Delete(models.Client{ID: id}).Error

	if err != nil {
		return err
	}

	return deleteUserById(client.UserID)
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

func GetUserById(id int64) (*models.User, error) {
	user := &models.User{}

	if err := db.First(&user, id).Error; err != nil {
		logrus.Info("USer not found: ", id)
		return nil, err
	}

	return user, nil
}

func createUser(email, password string, tx *gorm.DB) (*models.User, error) {
	passwordHash, err := utils.HashPassword(password)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         0,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return user, nil
}

func deleteUserById(id int64) error {
	return db.Delete(models.User{ID: id}).Error
}

func Feed(offset, limit int64) ([]models.Master, error) {
	var masters []models.Master
	if err := db.Limit(limit).Offset(offset).Find(&masters).Error; err != nil {
		return nil, err
	}
	return masters, nil
}
