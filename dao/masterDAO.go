package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// GetMasterByID returns master by id
func GetMasterByID(id int64) (*models.Master, error) {
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

// CreateMaster creates new master
func CreateMaster(username, email, password, name string, addressID, photoID int64) (*models.Master, error) {
	tx := db.Begin()

	user, err := createUser(email, password, tx)
	if err != nil {
		return nil, err
	}

	// TODO check addressID & photoID

	master := &models.Master{
		UserID:    user.ID,
		Name:      name,
		AddressID: addressID,
		PhotoID:   photoID,
	}

	if err := tx.Create(master).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return master, nil
}

// DeleteMaster deletes master
func DeleteMaster(id int64) error {
	master, err := GetMasterByID(id)

	if err != nil {
		return err
	}

	err = db.Delete(models.Master{ID: id}).Error

	if err != nil {
		return err
	}

	return deleteUser(master.UserID)
}

// MasterCount return count of masters
func MasterCount() int64 {
	var count int64
	if err := db.Model(&models.Master{}).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

// Feed returns feed
func Feed(offset, limit int) ([]*models.Master, error) {
	var masters []*models.Master
	if err := db.Limit(limit).Offset(offset).Find(&masters).Error; err != nil {
		return nil, err
	}
	return masters, nil
}
