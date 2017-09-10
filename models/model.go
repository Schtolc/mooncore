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

type Address struct {
	ID  int
	Lat float32
	Lon float32
}

type Photo struct {
	ID   int
	Path string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Address  Address
	Photo    Photo
}

type Master struct {
	ID                 int
	Name               string
	Address            Address
	Photo              Photo
	WorkingPlacePhotos []Photo `gorm:"many2many:working_place_photos;"`
}

type Service struct {
	ID                int
	Name              string
	Master            Master
	Descr             string // TODO make unlimited
	ManicureType      ManicureType
	ManicureMaterials []ManicureMaterial `gorm:"many2many:service_manicure_materials;"`
	Photos            []Photo            `gorm:"many2many:service_photos;"`
}

type ManicureType struct {
	ID   int
	Name string
}

type ManicureMaterial struct {
	ID       int
	Firm     string
	Palette  string
	Descr    string // TODO make unlimited
	Services []Service `gorm:"many2many:service_manicure_materials;"`
}
