package models

import (
	"time"
)

// Mock is a simple database model for checking if database can properly save and migrate models.
type Mock struct {
	ID   int    `gorm:"primary_key"`
	Path string `sql:"size:255;index"`
	Time time.Time
}
