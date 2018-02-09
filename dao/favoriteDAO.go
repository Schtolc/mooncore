package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// ClientFavorites returns client favorites
func ClientFavorites(client *models.Client) ([]*models.Master, error) {
	var masters []*models.Master

	if err := db.Model(client).Association("favorites").Find(&masters).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	return masters, nil
}
