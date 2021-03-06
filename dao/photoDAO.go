package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
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
func CreatePhoto(path string, tags []int64) (*models.Photo, error) {
	photo := &models.Photo{
		Path: path,
		Tags: []models.Tag{},
	}

	if err := db.Where("id in (?)", tags).Find(&photo.Tags).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	if err := db.Create(photo).Error; err != nil {
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
