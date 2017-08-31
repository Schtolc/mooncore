package models

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/google/uuid"
)


type Metric struct {
	Id   string
	Path string 	`sql:"size:255;index"`
	Time time.Time
}


func (m *Metric) BeforeCreate(scope *gorm.Scope) {
	if m.Id == "" {
		uuid := uuid.New().String()
		scope.SetColumn("id", uuid)
	}
}

