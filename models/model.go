package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Metric struct {
	Id   string
	Path string `sql:"size:255;index"`
	Time time.Time
}

func (m *Metric) BeforeCreate(scope *gorm.Scope) {
	if m.Id == "" {
		uuid := uuid.New().String()
		scope.SetColumn("id", uuid)
	}
}
