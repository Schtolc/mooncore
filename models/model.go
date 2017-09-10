package models

import (
	"time"
)

// Metric is a simple database model for checking if database can properly save and migrate models.
type Metric struct {
	ID   int    `gorm:"primary_key"`
	Path string `sql:"size:255;index"`
	Time time.Time
}

// Address model
type Address struct {
	ID  int
	Lat float32
	Lon float32
}

// Photo model
type Photo struct {
	ID   int
	Path string
}

// User model
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Address  Address
	Photo    Photo
}

// Master model
type Master struct {
	ID                 int
	Name               string
	Address            Address
	Photo              Photo
	Service            []Service
	WorkingPlacePhotos []Photo `gorm:"many2many:working_place_photos;"`
}

// Service model
type Service struct {
	ID                int
	MasterID          int
	Name              string
	Description       string `sql:"type:text"`
	ManicureType      ManicureType
	ManicureMaterials []ManicureMaterial `gorm:"many2many:service_manicure_materials;"`
	Photos            []Photo            `gorm:"many2many:service_photos;"`
}

// ManicureType model
type ManicureType struct {
	ID   int
	Name string
}

// ManicureMaterial model
type ManicureMaterial struct {
	ID          int
	Firm        string
	Palette     string
	Description string    `sql:"type:text"`
	Services    []Service `gorm:"many2many:service_manicure_materials;"`
}
