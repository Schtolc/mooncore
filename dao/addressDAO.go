package dao

import (
	"fmt"
	"github.com/Schtolc/mooncore/models"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
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

	line, station, err := getMetro(lat, lon)
	if err != nil {
		line = "Unknown"
		station = "Unknown"
	}

	color, err := getLineColor(line)
	if err != nil {
		color = "000000"
	}

	address := &models.Address{
		Lat:         lat,
		Lon:         lon,
		Description: description,
		Line:        line,
		Station:     station,
		Color:       color,
	}
	if err := db.Create(address).Error; err != nil {
		logrus.Error(err)
		return nil, err
	}
	return address, nil
}

// DeleteAddress deletes address
func DeleteAddress(id int64) error {
	if err := db.Delete(models.Address{ID: id}).Error; err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func getMetro(lat, lon float64) (string, string, error) {
	resp, err := http.Get(fmt.Sprintf("https://geocode-maps.yandex.ru/1.x/?format=json&geocode=%f,%f&kind=metro&results=1", lon, lat))
	if err != nil {
		return "", "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	details, _, _, err := jsonparser.Get(body, "response", "GeoObjectCollection", "featureMember", "[0]", "GeoObject", "metaDataProperty", "GeocoderMetaData", "AddressDetails", "Country", "AdministrativeArea", "Locality", "Thoroughfare")
	if err != nil {
		return "", "", err
	}

	line, err := jsonparser.GetUnsafeString(details, "ThoroughfareName")
	if err != nil {
		return "", "", err
	}

	station, err := jsonparser.GetUnsafeString(details, "Premise", "PremiseName")
	if err != nil {
		return "", "", err
	}

	return line, station, nil
}

func getLineColor(line string) (string, error) {
	lineColor := models.LineColor{}

	if err := db.Where("line = ?", line).First(&lineColor).Error; err != nil {
		logrus.Error("Unexpected line name: ", line)
		return "", err
	}

	return lineColor.Color, nil
}
