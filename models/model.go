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
