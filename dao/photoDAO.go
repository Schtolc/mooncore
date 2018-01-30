package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

func GetPhotoById(id int64) (*models.Photo, error) {
	photo := models.Photo{}

	if err := db.First(&photo, id).Error; err != nil {
		logrus.Error(err)
		return nil, nil
	}

	return &photo, nil
}

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
		return nil, err
	}

	return photo, nil
}

func DeletePhoto(id int64) error {
	return db.Delete(models.Photo{ID: id}).Error
}
