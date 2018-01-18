package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

func GetTagById(id int64) (*models.Tag, error) {
	tag := models.Tag{}
	if err := db.First(&tag, id).Error; err != nil {
		logrus.Error(err)
		return nil, nil
	}
	return &tag, nil
}

func CreateTag(name string) (*models.Tag, error) {
	tag := &models.Tag{
		Name: name,
	}

	if err := db.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}
