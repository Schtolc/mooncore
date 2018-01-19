package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

func GetAddressById(id int) (*models.Address, error) {
	address := models.Address{}
	if err := db.First(&address, id).Error; err != nil {
		logrus.Error(err)
		return nil, nil
	}
	return &address, nil
}

func GetAddressesInArea(lat1, lon1, lat2, lon2 float64) ([]models.Address, error) {
	var addresses []models.Address
	query := "lat > ? AND lat < ? AND lon < ? AND lon > ?"
	if err := db.Where(query, lat1, lat2, lon1, lon2).Find(&addresses).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return addresses, nil
}

func CreateAddress(lat, lon float64, description string) (*models.Address, error) {
	address := &models.Address{
		Lat:         lat,
		Lon:         lon,
		Description: description,
	}
	if err := db.Create(address).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return address, nil
}
