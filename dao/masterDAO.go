package dao

import (
	"github.com/Schtolc/mooncore/models"
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

func DeleteMaster(id int64) error {
	master, err := GetMasterById(id)

	if err != nil {
		return err
	}

	err = db.Delete(models.Master{ID: id}).Error

	if err != nil {
		return err
	}

	return deleteUser(master.UserID)
}
