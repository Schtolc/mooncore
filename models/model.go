package models

import (
	"time"
)

type Metric struct {
	Id   uint   `gorm:"primary_key"`
	Path string `sql:"size:255;index"`
	Time time.Time
}
