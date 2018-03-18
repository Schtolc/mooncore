package dao

import (
	"fmt"
	"github.com/Schtolc/mooncore/models"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
)

// GetAddressByID returns address by id
func GetAddressByID(id int64) (*models.Address, error) {
	address := models.Address{}
	if err := db.Preload("Stations").First(&address, id).Error; err != nil {
		logrus.Error(err)
		return nil, err
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
func CreateAddress(lat, lon float64, tx *gorm.DB) (*models.Address, error) {
	house, err := getHouse(lat, lon)
	if err != nil {
		logrus.Error(err)
		house = "Not found"
	}
	address := &models.Address{
		Lat:         lat,
		Lon:         lon,
		Description: house,
	}
	if err := tx.Create(address).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	metro, err := getMetro(lat, lon)
	for _, m := range metro {
		m.AddressID = address.ID

		if err := tx.Create(m).Error; err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return address, nil
}
// UpdateAddress old address
func UpdateAddress(id int64, lat, lon float64, tx *gorm.DB) error {
	house, err := getHouse(lat, lon)
	if err != nil {
		logrus.Error(err)
		house = "Not found"
	}
	address := &models.Address{
		ID:          id,
		Lat:         lat,
		Lon:         lon,
		Description: house,
	}
	if err := tx.Model(&models.Address{}).Update(address).Error; err != nil {
		logrus.Error(err)
		return err
	}
	metro, err := getMetro(lat, lon)
	for _, m := range metro {
		m.AddressID = address.ID
		if err := tx.Model(&models.AddressMetro{}).Update(m).Error; err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

// DeleteAddress deletes address
func DeleteAddress(id int64) error {
	if err := db.Delete(models.Address{ID: id}).Error; err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func getHouse(lat, lon float64) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://geocode-maps.yandex.ru/1.x/?format=json&geocode=%f,%f&kind=house&results=1", lon, lat))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	details, _, _, err := jsonparser.Get(body, "response", "GeoObjectCollection", "featureMember", "[0]", "GeoObject", "metaDataProperty", "GeocoderMetaData", "AddressDetails", "Country", "AdministrativeArea", "Locality", "Thoroughfare")
	if err != nil {
		return "", err
	}

	street, err := jsonparser.GetUnsafeString(details, "ThoroughfareName")
	if err != nil {
		return "", err
	}

	house, err := jsonparser.GetUnsafeString(details, "Premise", "PremiseNumber")
	if err != nil {
		return "", err
	}

	return street + " " + house, nil
}

func getMetro(lat, lon float64) ([]*models.AddressMetro, error) {
	var stations []*models.MetroStation

	if err := db.Preload("Line").Where("pow(? - lat, 2) + pow(cos(pi() * ? / 180) * (? - lon), 2) < pow(? / 100, 2)", lat, lat, lon, 2).Find(&stations).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}

	addresses := make([]*models.AddressMetro, len(stations))
	for i, station := range stations {
		addresses[i] = &models.AddressMetro{
			Name:     station.Name,
			Line:     station.Line.Name,
			Color:    station.Line.Color,
			Distance: math.Sqrt(math.Pow(station.Lat-lat, 2)+math.Pow(math.Cos(math.Pi*station.Lat/180)*(station.Lon-lon), 2)) * 100,
		}
	}

	sort.Slice(addresses, func(i, j int) bool {
		return addresses[i].Distance < addresses[j].Distance
	})

	return addresses, nil
}
