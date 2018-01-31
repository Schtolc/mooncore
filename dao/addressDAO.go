package dao

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/sirupsen/logrus"
)

// GetAddressByID returns address by id
func GetAddressByID(id int64) (*models.Address, error) {
	address := models.Address{}
	if err := db.First(&address, id).Error; err != nil {
		logrus.Error(err)
		return nil, nil
	}
	return &address, nil
}

// GetAddressesInArea returns addresses in area
func GetAddressesInArea(lat1, lon1, lat2, lon2 float64) ([]models.Address, error) {
	var addresses []models.Address
	query := "lat > ? AND lat < ? AND lon < ? AND lon > ?"
	if err := db.Where(query, lat1, lat2, lon1, lon2).Find(&addresses).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return addresses, nil
}

// CreateAddress creates new address
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

// DeleteAddress deletes address
func DeleteAddress(id int64) error {
	return db.Delete(models.Address{ID: id}).Error
}
