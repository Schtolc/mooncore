package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// GetTagByID returns tag by id
func GetTagByID(id int64) (*models.Tag, error) {
	tag := models.Tag{}
	if err := db.First(&tag, id).Error; err != nil {
		logrus.Error(err)
		return nil, nil
	}
	return &tag, nil
}

// CreateTag creates tag
func CreateTag(name string) (*models.Tag, error) {
	tag := &models.Tag{
		Name: name,
	}

	if err := db.Create(tag).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	return tag, nil
}

// PhotoTags returns photo tag
func PhotoTags(photo *models.Photo) ([]*models.Tag, error) {
	var tags []*models.Tag

	if err := db.Model(photo).Association("tags").Find(&tags).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	return tags, nil
}
