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

// Address is a model for addresses
type Address struct {
	ID  int
	Lat float32
	Lon float32
}

// Photo is a model for paths to photos
type Photo struct {
	ID   int
	Path string
}

// User is a model for store users
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Address  Address
	Photo    Photo
}

// Master is a model for store masters
type Master struct {
	ID                 int
	Name               string
	Address            Address
	Photo              Photo
	WorkingPlacePhotos []Photo `gorm:"many2many:working_place_photos;"`
}

// Service is a model for store services
type Service struct {
	ID                int
	Name              string
	Master            Master
	Descr             string // TODO make unlimited
	ManicureType      ManicureType
	ManicureMaterials []ManicureMaterial `gorm:"many2many:service_manicure_materials;"`
	Photos            []Photo            `gorm:"many2many:service_photos;"`
}

// ManicureType is a model for store types of manicure
type ManicureType struct {
	ID   int
	Name string
}

// ManicureMaterial is a model for store types of materials
type ManicureMaterial struct {
	ID       int
	Firm     string
	Palette  string
	Descr    string    // TODO make unlimited
	Services []Service `gorm:"many2many:service_manicure_materials;"`
}
