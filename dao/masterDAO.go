package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"errors"
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

func GetMasterFromContext(p graphql.ResolveParams)(*models.Master, error) {
	user := p.Context.Value(utils.GraphQLContextUserKey)
	if user == nil {
		return nil, errors.New("no user")
	}
	userModel := user.(*models.User)
	master := &models.Master{}
	if err := db.Where("user_id=?",userModel.ID).First(master).Error; err != nil {
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
		UserID:    user.ID,
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

func UpdateMaster(master *models.Master, tx *gorm.DB) error {
	if tx == nil { tx = db }
	if err := tx.Model(&models.Master{}).Updates(master).Error; err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}