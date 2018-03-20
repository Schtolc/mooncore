package dao

import (
	"database/sql"
	"errors"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

// GetMasterByID returns master by id
func GetMasterByID(id int64) (*models.Master, error) {
	master := &models.Master{}
	if err := db.Preload("User").First(master, id).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return master, nil
}

// GetMasterFromContext get master from context
func GetMasterFromContext(p graphql.ResolveParams) (*models.Master, error) {
	user := p.Context.Value(utils.GraphQLContextUserKey)
	if user == nil {
		return nil, errors.New("no user")
	}
	userModel := user.(*models.User)
	master := &models.Master{}
	if err := db.Where("user_id=?", userModel.ID).First(master).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return master, nil
}

// CreateMaster creates new master
func CreateMaster(email, passwordHash string) (*models.Master, error) {
	tx := db.Begin()
	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         models.MasterRole,
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}
	master := &models.Master{
		UserID: user.ID,
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
	tx := db.Begin()
	if err = tx.Delete(models.Master{ID: master.ID}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}

	if err := tx.Delete(models.User{ID: master.UserID}).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return err
	}
	tx.Commit()
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

// MasterSigns returns master signs
func MasterSigns(master *models.Master) ([]*models.Sign, error) {
	var signs []*models.Sign

	if err := db.Model(master).Association("signs").Find(&signs).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	return signs, nil
}

// EditMaster update master parameters
func EditMaster(master *models.Master, name, photoPath string, lat, lon float64) (*models.Master, error) {
	tx := db.Begin()
	photo := &models.Photo{
		Path: photoPath,
		Tags: []models.Tag{},
	}
	if master.PhotoID.Valid {
		photo.ID = master.PhotoID.Int64
		if err := tx.Model(photo).Update(photo).Error; err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}
	} else {
		if err := tx.Create(photo).Error; err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}
		master.PhotoID = sql.NullInt64{
			Int64: photo.ID,
			Valid: true,
		}
	}

	house, err := getHouse(lat, lon)
	if err != nil {
		logrus.Error(err)
		house = "Not found"
	}
	address := &models.Address{
		Lat:         lat,
		Lon:         lon,
		Description: house,
	}
	if master.AddressID.Valid {
		address.ID = master.AddressID.Int64
		if err := tx.Model(&models.Address{}).Update(address).Error; err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}
		if err := tx.Where("address_id=?", address.ID).Delete(models.AddressMetro{}).Error; err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}
	} else {
		if err := tx.Create(address).Error; err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}
		master.AddressID = sql.NullInt64{
			Int64: address.ID,
			Valid: true,
		}
	}
	metroArr, err := getMetro(lat, lon)
	for _, metro := range metroArr {
		metro.AddressID = address.ID
		if err := tx.Create(metro).Error; err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}
	}
	master.Name = name
	if err := tx.Model(&models.Master{}).Updates(master).Error; err != nil {
		tx.Rollback()
		logrus.Error(err)
		return nil, err
	}
	tx.Commit()
	return master, nil
}
