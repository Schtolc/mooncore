package models

import (
	"time"
)

type Metric struct {
	Id   int    `gorm:"primary_key"`
	Path string `sql:"size:255;index"`
	Time time.Time
}
