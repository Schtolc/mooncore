package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// MasterServices returns master services
func MasterServices(master *models.Master) ([]*models.Service, error) {
	var services []*models.Service

	if err := db.Where("master_id = ?", master.ID).Find(&services).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	return services, nil
}
