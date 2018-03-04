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
func CreateMaster(username, email, password, name string, addressID int64) (*models.Master, error) {
	tx := db.Begin()

	user, err := createUser(email, password, models.MasterRole, tx)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// TODO check addressID & photoID

	master := &models.Master{
		UserID:    user.ID,
		Name:      name,
		AddressID: addressID,
	}

	if err := tx.Create(master).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}

	tx.Commit()
	return master, nil
}

// DeleteMaster deletes master
func DeleteMaster(id int64) error {
	master, err := GetMasterByID(id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = db.Delete(models.Master{ID: id}).Error
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = deleteUser(master.UserID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

// MasterCount return count of masters
func MasterCount() (int64, error) {
	var count int64
	if err := db.Model(&models.Master{}).Count(&count).Error; err != nil {
		logrus.Error(err)
		return 0, err
	}
	return count, nil
}

// Feed returns feed
func Feed(offset, limit int) ([]*models.Master, error) {
	var masters []*models.Master
	if err := db.Limit(limit).Offset(offset).Preload("Photos").Find(&masters).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return masters, nil
}

// MasterSigns returns master signs
func MasterSigns(master *models.Master) ([]*models.Sign, error) {
	var signs []*models.Sign

	if err := db.Model(master).Association("signs").Find(&signs).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	return signs, nil
}
