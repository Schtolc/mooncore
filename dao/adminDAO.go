package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// GetAdminByID returns Admin by id
func GetAdminByID(id int64) (*models.Admin, error) {
	user := &models.User{}
	admin := &models.Admin{}
	if dbc := db.First(admin, id); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, nil
	}
	if dbc := db.First(user, admin.UserID); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	admin.User = *user
	return admin, nil
}

// CreateAdmin creates new Admin
func CreateAdmin(email, passwordHash string) (*models.Admin, error) {
	tx := db.Begin()
	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         models.AdminRole,
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}
	admin := &models.Admin{
		UserID: user.ID,
	}
	if err := tx.Create(admin).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}
	tx.Commit()
	return admin, nil
}

// DeleteAdmin deletes Admin
func DeleteAdmin(id int64) error {
	admin, err := GetAdminByID(id)
	if err != nil {
		logrus.Error(err)
		return err
	}
	tx := db.Begin()
	if err = tx.Delete(models.Admin{ID: admin.ID}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}

	if err := tx.Delete(models.User{ID: admin.UserID}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}
	tx.Commit()
	return nil
}
