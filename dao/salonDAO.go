package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// GetSalonByID returns Salon by id
func GetSalonByID(id int64) (*models.Salon, error) {
	user := &models.User{}
	salon := &models.Salon{}
	if dbc := db.First(salon, id); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, nil
	}
	if dbc := db.First(user, salon.UserID); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	salon.User = *user
	return salon, nil
}

// CreateSalon creates new salon
func CreateSalon(email, passwordHash string) (*models.Salon, error) {
	tx := db.Begin()
	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         models.SalonRole,
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}
	salon := &models.Salon{
		UserID: user.ID,
	}
	if err := tx.Create(salon).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}
	tx.Commit()
	return salon, nil
}

// DeleteSalon deletes Salon
func DeleteSalon(id int64) error {
	salon, err := GetSalonByID(id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	tx := db.Begin()
	if err = tx.Delete(models.Salon{ID: salon.ID}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}

	if err := tx.Delete(models.User{ID: salon.UserID}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}
	tx.Commit()
	return nil
}
