package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// Feed returns feed
func Feed(offset, limit int) ([]*models.Master, error) {
	var masters []*models.Master
	if err := db.Limit(limit).Offset(offset).
		Where("photo_id is not null AND address_id is not null").
		Preload("Photos").Preload("Signs").Find(&masters).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return masters, nil
}
