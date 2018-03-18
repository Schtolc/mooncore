package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

// GetPhotoByID returns photo by id
func GetPhotoByID(id int64) (*models.Photo, error) {
	photo := models.Photo{}

	if err := db.First(&photo, id).Error; err != nil {
		logrus.Error(err)
		return nil, nil
	}

	return &photo, nil
}

// CreatePhoto creates new photo
func CreatePhoto(path string, tags []int64, tx *gorm.DB) (*models.Photo, error) {
	if tx == nil { tx = db }
	photo := &models.Photo{
		Path: path,
		Tags: []models.Tag{},
	}

	if err := tx.Where("id in (?)", tags).Find(&photo.Tags).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	if err := tx.Create(photo).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	return photo, nil
}

// DeletePhoto deletes photo
func DeletePhoto(id int64) error {
	if err := db.Delete(models.Photo{ID: id}).Error; err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// MasterPhotos returns master photos
func MasterPhotos(master *models.Master) ([]*models.Photo, error) {
	var photos []*models.Photo

	if err := db.Model(master).Association("photos").Find(&photos).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	return photos, nil
}

// UpdatePhoto  update photo
func UpdatePhoto(id int64, path string, tags []int64, tx *gorm.DB) error {
	if tx == nil { tx = db }
	photo := &models.Photo{
		ID: id,
		Path: path,
		Tags: []models.Tag{},
	}
	if err := tx.Where("id in (?)", tags).Find(&photo.Tags).Error; err != nil {
		logrus.Error(err)
		return err
	}
	if err := tx.Model(photo).Update(photo).Error; err != nil {
		return err
	}
	return nil;
}